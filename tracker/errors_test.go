package tracker

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestErrorResponse_Error(t *testing.T) {
	errResp := &ErrorResponse{
		Response: &http.Response{
			StatusCode: 404,
			Request: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Scheme: "https",
					Host:   "api.tracker.yandex.net",
					Path:   "/v2/issues/QUEUE-1",
				},
			},
		},
		ErrorMessages: []string{"Issue not found"},
		Errors:        map[string]string{"key": "not found"},
	}

	msg := errResp.Error()

	if !strings.Contains(msg, "GET") {
		t.Errorf("Error() = %q, want to contain method GET", msg)
	}
	if !strings.Contains(msg, "/v2/issues/QUEUE-1") {
		t.Errorf("Error() = %q, want to contain URL path", msg)
	}
	if !strings.Contains(msg, "404") {
		t.Errorf("Error() = %q, want to contain status code 404", msg)
	}
	if !strings.Contains(msg, "Issue not found") {
		t.Errorf("Error() = %q, want to contain error message", msg)
	}
	if !strings.Contains(msg, "not found") {
		t.Errorf("Error() = %q, want to contain errors map value", msg)
	}
}

func TestIsNotFound(t *testing.T) {
	t.Run("404 ErrorResponse", func(t *testing.T) {
		err := &ErrorResponse{
			Response: &http.Response{StatusCode: http.StatusNotFound},
		}
		if !IsNotFound(err) {
			t.Error("IsNotFound returned false for 404, want true")
		}
	})

	t.Run("200 ErrorResponse", func(t *testing.T) {
		err := &ErrorResponse{
			Response: &http.Response{StatusCode: http.StatusOK},
		}
		if IsNotFound(err) {
			t.Error("IsNotFound returned true for 200, want false")
		}
	})

	t.Run("non-ErrorResponse error", func(t *testing.T) {
		err := fmt.Errorf("some other error")
		if IsNotFound(err) {
			t.Error("IsNotFound returned true for non-ErrorResponse, want false")
		}
	})
}

func TestIsForbidden(t *testing.T) {
	t.Run("403 ErrorResponse", func(t *testing.T) {
		err := &ErrorResponse{
			Response: &http.Response{StatusCode: http.StatusForbidden},
		}
		if !IsForbidden(err) {
			t.Error("IsForbidden returned false for 403, want true")
		}
	})

	t.Run("200 ErrorResponse", func(t *testing.T) {
		err := &ErrorResponse{
			Response: &http.Response{StatusCode: http.StatusOK},
		}
		if IsForbidden(err) {
			t.Error("IsForbidden returned true for 200, want false")
		}
	})

	t.Run("non-ErrorResponse error", func(t *testing.T) {
		err := fmt.Errorf("some other error")
		if IsForbidden(err) {
			t.Error("IsForbidden returned true for non-ErrorResponse, want false")
		}
	})
}

func TestIsRateLimited(t *testing.T) {
	t.Run("429 ErrorResponse", func(t *testing.T) {
		err := &ErrorResponse{
			Response: &http.Response{StatusCode: http.StatusTooManyRequests},
		}
		if !IsRateLimited(err) {
			t.Error("IsRateLimited returned false for 429, want true")
		}
	})

	t.Run("200 ErrorResponse", func(t *testing.T) {
		err := &ErrorResponse{
			Response: &http.Response{StatusCode: http.StatusOK},
		}
		if IsRateLimited(err) {
			t.Error("IsRateLimited returned true for 200, want false")
		}
	})

	t.Run("non-ErrorResponse error", func(t *testing.T) {
		err := fmt.Errorf("some other error")
		if IsRateLimited(err) {
			t.Error("IsRateLimited returned true for non-ErrorResponse, want false")
		}
	})
}

// Suppress unused import warning.
var _ = errors.As
