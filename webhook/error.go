package webhook

import "fmt"

// RequestError provides informations about generic request errors.
type RequestError struct {
	HTTPStatusCode int
	Err            error
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("error, status code: %d, message: %s", e.HTTPStatusCode, e.Err)
}

func (e *RequestError) Unwrap() error {
	return e.Err
}
