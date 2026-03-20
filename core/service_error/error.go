// -----------------------------------------------------------------------------
// Copyright (C) 2026 tadnavel
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// -----------------------------------------------------------------------------

package svce

import (
	"errors"
	"fmt"
	"net/http"
)

// define error catalog:
// ErrUserNotFound ErrCode = "U001"
// ErrUserALreadyExist ErrCode = "U002"
// ....
type ErrCode string

type ServiceError struct {
	Code       ErrCode                `json:"code"`
	HTTPStatus int                    `json:"-"`
	Message    string                 `json:"message"`
	Fields     map[string]any `json:"fields,omitempty"`

	rootError error
	logMsg    string
}

type Option func(*ServiceError)

func WithField(key string, value any) Option {
	return func(e *ServiceError) {
		if e.Fields == nil {
			e.Fields = make(map[string]any)
		}
		e.Fields[key] = value
	}
}

func WithLogMessage(msg string) Option {
	return func(e *ServiceError) {
		e.logMsg = msg
	}
}

func NewServiceError(status int, code ErrCode, message string, rootError error, opts ...Option) *ServiceError {
	e := &ServiceError{
		Code:       code,
		HTTPStatus: status,
		Message:    message,
		rootError:  rootError,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *ServiceError) Error() string {
	if e.rootError != nil {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.rootError)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// use error.As to get service error from error
func (e *ServiceError) Unwrap() error {
	return e.rootError
}

func AsServiceError(err error) (*ServiceError, bool) {
	var sErr *ServiceError
	return sErr, errors.As(err, &sErr)
}

func (e *ServiceError) GetFields() map[string]any {
	fields := map[string]any{
		"error_code":  string(e.Code),
		"http_status": e.HTTPStatus,
	}
	if e.rootError != nil {
		fields["root_error"] = e.rootError.Error()
	}
	for k, v := range e.Fields {
		fields[k] = v
	}
	return fields
}

// 400
func NewBadRequest(code ErrCode, message string, rootError error, opts ...Option) *ServiceError {
	return NewServiceError(http.StatusBadRequest, code, message, rootError, opts...)
}

// 401
func NewUnauthorized(code ErrCode, message string, rootError error, opts ...Option) *ServiceError {
	return NewServiceError(http.StatusUnauthorized, code, message, rootError, opts...)
}

// 403
func NewForbidden(code ErrCode, message string, rootError error, opts ...Option) *ServiceError {
	return NewServiceError(http.StatusForbidden, code, message, rootError, opts...)
}

// 404
func NewNotFound(code ErrCode, message string, rootError error, opts ...Option) *ServiceError {
	return NewServiceError(http.StatusNotFound, code, message, rootError, opts...)
}

// 408
func NewRequestTimeout(code ErrCode, message string, rootError error, opts ...Option) *ServiceError {
	return NewServiceError(http.StatusRequestTimeout, code, message, rootError, opts...)
}

// 409
func NewConflict(code ErrCode, message string, rootError error, opts ...Option) *ServiceError {
	return NewServiceError(http.StatusConflict, code, message, rootError, opts...)
}

// 429
func NewTooManyRequests(code ErrCode, message string, rootError error, opts ...Option) *ServiceError {
	return NewServiceError(http.StatusTooManyRequests, code, message, rootError, opts...)
}

// 500
func NewInternalServerError(code ErrCode, message string, rootError error, opts ...Option) *ServiceError {
	return NewServiceError(http.StatusInternalServerError, code, message, rootError, opts...)
}

// 502
func NewBadGateway(code ErrCode, message string, rootError error, opts ...Option) *ServiceError {
	return NewServiceError(http.StatusBadGateway, code, message, rootError, opts...)
}

// 503
func NewServiceUnavailable(code ErrCode, message string, rootError error, opts ...Option) *ServiceError {
	return NewServiceError(http.StatusServiceUnavailable, code, message, rootError, opts...)
}

// 504
func NewGatewayTimeout(code ErrCode, message string, rootError error, opts ...Option) *ServiceError {
	return NewServiceError(http.StatusGatewayTimeout, code, message, rootError, opts...)
}
