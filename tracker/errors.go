package tracker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ErrorResponse represents an error response from the Yandex Tracker API.
// It implements the error interface.
type ErrorResponse struct {
	// Response is the HTTP response that caused this error.
	Response *http.Response `json:"-"`

	// Errors is a map of field-level errors returned by the API.
	Errors map[string]string `json:"errors"`

	// ErrorMessages is a list of error messages returned by the API.
	ErrorMessages []string `json:"errorMessages"`
}

// Error implements the error interface.
// It includes the HTTP method, URL, status code, error messages, and errors map.
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %v",
		r.Response.Request.Method,
		r.Response.Request.URL,
		r.Response.StatusCode,
		r.ErrorMessages,
		r.Errors,
	)
}

// checkResponse checks the API response for errors. If the response
// indicates an error (status code outside 2xx range), it returns an
// ErrorResponse with the parsed error body.
func checkResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c < 300 {
		return nil
	}

	errResp := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, errResp) // best-effort decode of error body
	}

	return errResp
}

// IsNotFound returns true if the error is an ErrorResponse with a 404 status code.
func IsNotFound(err error) bool {
	var errResp *ErrorResponse
	if errors.As(err, &errResp) {
		return errResp.Response.StatusCode == http.StatusNotFound
	}
	return false
}

// IsForbidden returns true if the error is an ErrorResponse with a 403 status code.
func IsForbidden(err error) bool {
	var errResp *ErrorResponse
	if errors.As(err, &errResp) {
		return errResp.Response.StatusCode == http.StatusForbidden
	}
	return false
}

// IsRateLimited returns true if the error is an ErrorResponse with a 429 status code.
func IsRateLimited(err error) bool {
	var errResp *ErrorResponse
	if errors.As(err, &errResp) {
		return errResp.Response.StatusCode == http.StatusTooManyRequests
	}
	return false
}
