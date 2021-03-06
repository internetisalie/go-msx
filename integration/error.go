// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package integration

import (
	"github.com/pkg/errors"
	"net/http"
)

// NOTE: API response error payloads are in response.go

type ResponseError interface {
	Error() error
	IsError() bool
}

type StatusError struct {
	Code int
	Body string
	Err  error
}

// Allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Err.Error()
}

func (se StatusError) StatusCode() int {
	return se.Code
}

func NewStatusError(statusCode int, body []byte) error {
	return &StatusError{
		Code: statusCode,
		Err:  errors.New(http.StatusText(statusCode)),
		Body: string(body),
	}
}
