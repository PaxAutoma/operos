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

package widgets

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
	Show    bool
}

func NewValidationError(field string, message string) *ValidationError {
	return &ValidationError{field, message, true}
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", ve.Field, ve.Message)
}

func JoinValidationErrors(errors []error) string {
	errStrings := make([]string, 0)
	for _, err := range errors {
		if ve, ok := err.(*ValidationError); ok && ve.Show {
			errStrings = append(errStrings, err.Error())
		}
	}

	return strings.Join(errStrings, "\n")
}

func ValidateNotEmpty(field string, value string) []error {
	value = strings.TrimSpace(value)

	if value == "" {
		return []error{
			NewValidationError(field, "cannot be empty"),
		}
	}

	return []error{}
}

func ValidateIP(field string, value string) []error {
	if net.ParseIP(value) == nil {
		return []error{
			NewValidationError(field, "not a valid IP"),
		}
	}
	return []error{}
}

func ValidateIPNet(field string, value string) []error {
	if _, _, err := net.ParseCIDR(value); err != nil {
		return []error{
			NewValidationError(field, "not a valid CIDR string"),
		}
	}
	return []error{}
}

func ValidateIntMinMax(field string, value string, min int, max int) []error {
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return []error{
			NewValidationError(field, "is not a valid integer"),
		}
	}

	if intVal < min || intVal > max {
		return []error{
			NewValidationError(field, "is out of bounds"),
		}
	}
	return []error{}
}
