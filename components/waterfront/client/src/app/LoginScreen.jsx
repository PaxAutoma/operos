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

import React from 'react';
import PropTypes from 'prop-types';
import {withStyles} from 'material-ui/styles';
import Card, {CardContent, CardMedia} from 'material-ui/Card';
import TextField from 'material-ui/TextField';
import Button from 'material-ui/Button';
import {CircularProgress} from 'material-ui/Progress';
import {graphql, gql} from 'react-apollo';
import {red} from 'material-ui/colors';

const styles = theme => ({
  root: {
    width: '100%',
    height: '100%',
    display: 'flex',
    justifyContent: 'center',
    background: 'radial-gradient(circle at center top 180px, white, gray)'
  },
  loginCard: {
    marginTop: 80,
    width: 400,
  },
  media: {
    height: 116
  },
  form: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center'
  },
  loginButton: {
    marginTop: 24,
  },
  textbox: {
    width: '90%',
    marginTop: theme.spacing.unit * 2
  },
  failed: {
    fontWeight: 'bold',
    color: red[500]
  }
});

class LoginScreen extends React.Component {
  constructor() {
    super()
    this.state = {
      loggingIn: false,
      username: '',
      password: '',
      failed: false
    }
  }

  onChange(event) {
    this.setState({
      [event.target.id]: event.target.value
    })
  }

  onLoginClick(evt) {
    evt.preventDefault();

    this.setState({
      failed: false,
      loggingIn: true
    })

    this.props.mutate({
      variables: {
        username: this.state.username,
        password: this.state.password
      }
    })
    .then(res => {
      if (!res.data.login.logged_in) {
        this.setState({
          loggingIn: false,
          failed: true
        });
      }
      // If we've logged in successfully then this screen will no longer be
      // visible so we don't have to set anything here.
    })
    .catch(err => {
      this.setState({
        loggingIn: false
      });
    });
  }

  render() {
    const {classes} = this.props;

    return (
      <div className={classes.root}>
        <div>
          <Card className={classes.loginCard} raised>
            <CardMedia image="/static/login-top.png" title="Operos Login" className={classes.media} />
            <CardContent>
              <form noValidate autoComplete="off" className={classes.form}>
                <TextField
                    id="username"
                    label="Username"
                    className={classes.textbox}
                    value={this.state.username}
                    onChange={this.onChange.bind(this)}
                    autoFocus
                />
                <TextField
                    id="password"
                    label="Password"
                    type="password"
                    className={classes.textbox}
                    value={this.state.password}
                    onChange={this.onChange.bind(this)}
                />
                { this.state.loggingIn
                  ? <CircularProgress className={classes.loginButton} size={36} />
                  : <Button
                        raised
                        color="primary"
                        className={classes.loginButton}
                        onClick={this.onLoginClick.bind(this)}
                        type="submit"
                    >
                      Log in
                    </Button>
                }
                { this.state.failed &&
                  <p className={classes.failed}>Invalid username or password</p> }
              </form>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }
}

export default withStyles(styles)(graphql(gql`
  mutation login($username: String!, $password: String!) {
    login(username: $username, password: $password) {
      logged_in
      user {
        username
      }
    }
  }
`)(LoginScreen));
