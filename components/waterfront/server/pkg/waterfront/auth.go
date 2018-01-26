/*
Copyright 2018 Pax Automa Systems, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package waterfront

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	gorilla_context "github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/msteinert/pam"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func init() {
	gob.Register(time.Time{})
	gob.Register(User{})
}

type User struct {
	Username  string    `json:"username"`
	LoginTime time.Time `json:"login_time"`
}

type ContextKey string

const (
	ContextKeyuser = ContextKey("user")
)

type auth struct {
	store         sessions.Store
	validity      time.Duration
	authenticator Authenticator
}

func AuthSessionMiddleware(opts ...AuthSessionOption) *auth {
	ash := &auth{
		validity:      24 * time.Hour,
		store:         sessions.NewCookieStore([]byte("this is a totally secret key")),
		authenticator: PAMAuthenticator,
	}

	for _, opt := range opts {
		opt(ash)
	}

	return ash
}

type AuthSessionOption func(*auth)

func AuthValidity(validity time.Duration) AuthSessionOption {
	return func(ash *auth) {
		ash.validity = validity
	}
}

func SessionStore(store sessions.Store) AuthSessionOption {
	return func(ash *auth) {
		ash.store = store
	}
}

func (h *auth) getUser(r *http.Request) (*sessions.Session, *User) {
	session, err := h.store.Get(r, "operos-waterfront")
	if err != nil {
		log.Errorf("failed to obtain session: %s", err.Error())
	}

	var user User

	if username, password, ok := r.BasicAuth(); ok {
		var loggedIn bool
		if loggedIn, err = h.authenticator.Authenticate(username, password); err != nil {
			panic(errors.Wrap(err, "failed to authenticate"))
		}

		if !loggedIn {
			return session, nil
		}

		user = User{
			Username:  username,
			LoginTime: time.Now(),
		}
	} else {
		user, ok := session.Values["user"].(User)
		if !ok || user.LoginTime.Add(h.validity).Before(time.Now()) {
			return session, nil
		}
	}

	return session, &user
}

func (h *auth) GetAuthHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, user := h.getUser(r)

		if user == nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// This is necessary to prevent https://github.com/gorilla/sessions/issues/80
		defer gorilla_context.Clear(r)

		ctx := context.WithValue(r.Context(), ContextKeyuser, user)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *auth) GetLoginHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, user := h.getUser(r)
		h.login(w, r, user, session)
	})
}

func (h *auth) GetLogoutHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.getUser(r)
		h.logout(w, r, session)
	})
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *auth) login(w http.ResponseWriter, r *http.Request, user *User, s *sessions.Session) {
	loggedIn := user != nil

	if r.Method == "POST" {
		var creds credentials
		var err error

		if err = json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %s", err.Error()), http.StatusBadRequest)
			return
		}

		if loggedIn, err = h.authenticator.Authenticate(creds.Username, creds.Password); err != nil {
			panic(errors.Wrap(err, "failed to authenticate"))
		}

		if loggedIn {
			log.Infof("user %s authenticated successfully", creds.Username)
			user = &User{
				Username:  creds.Username,
				LoginTime: time.Now(),
			}
			s.Values["user"] = user
			if err := s.Save(r, w); err != nil {
				panic(errors.Wrap(err, "failed to save session"))
			}
		} else {
			log.Infof("user %s authentication failed", creds.Username)
		}
	}

	var response interface{}

	if loggedIn {
		response = struct {
			LoggedIn bool  `json:"logged_in"`
			User     *User `json:"user"`
		}{
			LoggedIn: true,
			User:     user,
		}
	} else {
		response = struct {
			LoggedIn bool `json:"logged_in"`
		}{
			LoggedIn: false,
		}
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		panic(err)
	}
}

func (h *auth) logout(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	delete(s.Values, "user")
	if err := s.Save(r, w); err != nil {
		log.Errorf("failed to save session: %s", err.Error())
	}

	response := struct {
		LoggedIn bool `json:"logged_in"`
	}{
		LoggedIn: false,
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		panic(err)
	}
}

type Authenticator interface {
	Authenticate(username, password string) (bool, error)
}

type AuthenticatorFunc func(username, password string) (bool, error)

func (a AuthenticatorFunc) Authenticate(username, password string) (bool, error) {
	return a(username, password)
}

var PAMAuthenticator = AuthenticatorFunc(func(username, password string) (bool, error) {
	t, err := pam.StartFunc("", username, func(style pam.Style, msg string) (string, error) {
		switch style {
		case pam.PromptEchoOff:
			return password, nil
		case pam.PromptEchoOn, pam.ErrorMsg, pam.TextInfo:
			return "", nil
		}

		return "", errors.New("Unrecognized PAM message style")
	})
	if err != nil {
		return false, errors.Wrap(err, "error setting up PAM authenticator")
	}

	err = t.Authenticate(0)
	if err != nil {
		return false, nil
	}

	return true, nil
})
