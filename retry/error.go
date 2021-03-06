// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package retry

type TransientError struct {
	Cause error
}

func (e *TransientError) Unwrap() error {
	return e.Cause
}

func (e *TransientError) IsPermanent() bool {
	return false
}

func (e *TransientError) Error() string {
	return e.Cause.Error()
}

func TransientErrorInterceptor(fn Retryable) error {
	return TransientErrorDecorator(fn)()
}

// TransientErrorDecorator wraps the supplied retryable to force it to return a TransientError
func TransientErrorDecorator(fn Retryable) Retryable {
	return func() error {
		err := fn()
		if err != nil {
			return &TransientError{
				Cause: err,
			}
		}

		return nil
	}
}

type PermanentError struct {
	Cause error
}

func (e *PermanentError) Unwrap() error {
	return e.Cause
}

func (e *PermanentError) IsPermanent() bool {
	return true
}

func (e *PermanentError) Error() string {
	return e.Cause.Error()
}

func PermanentErrorInterceptor(fn Retryable) error {
	return PermanentErrorDecorator(fn)()
}

// TransientErrorDecorator wraps the supplied retryable to force it to return a PermanentError
func PermanentErrorDecorator(fn Retryable) Retryable {
	return func() error {
		err := fn()
		if err != nil {
			return &PermanentError{
				Cause: err,
			}
		}

		return nil
	}
}
