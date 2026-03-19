// -----------------------------------------------------------------------------
// Copyright (C) 2026 tadnavel
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
	Fields     map[string]interface{} `json:"fields,omitempty"`

	rootError error
	logMsg    string
}

type Option func(*ServiceError)

func WithField(key string, value interface{}) Option {
	return func(e *ServiceError) {
		if e.Fields == nil {
			e.Fields = make(map[string]interface{})
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
