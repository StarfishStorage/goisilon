package api

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func makeResp(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}

func TestParseJSONError_WellFormed404(t *testing.T) {
	r := makeResp(404, `{"errors":[{"code":"ENOENT","message":"not found"}]}`)
	err := parseJSONError(r)
	je, ok := err.(*JSONError)
	if !ok {
		t.Fatalf("want *JSONError, got %T", err)
	}
	if je.StatusCode != 404 {
		t.Fatalf("want 404, got %d", je.StatusCode)
	}
	if got := je.Error(); got != "not found" {
		t.Fatalf("unexpected message: %q", got)
	}
}

func TestParseJSONError_SuccessShapedBody404(t *testing.T) {
	r := makeResp(404, `{"id":1}`)
	err := parseJSONError(r)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(*JSONError); !ok {
		t.Fatalf("want *JSONError, got %T", err)
	}
}

func TestParseJSONError_EmptyBody404(t *testing.T) {
	r := makeResp(404, ``)
	err := parseJSONError(r)
	je := err.(*JSONError)
	if je.Error() == "" {
		t.Fatalf("expected non-empty message, got empty")
	}
}

func TestJSONError_Error_NoErrs(t *testing.T) {
	je := &JSONError{StatusCode: 400, Err: nil}
	if got := je.Error(); got == "" {
		t.Fatalf("want non-empty fallback message, got empty")
	}
}
