package mcp

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLimitMCPRequestBodyRejectsOversizedContentLength(t *testing.T) {
	called := false
	handler := limitMCPRequestBody(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		called = true
	}), 8)

	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader("123456789"))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if called {
		t.Fatal("next handler was called for oversized MCP body")
	}
	if recorder.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusRequestEntityTooLarge)
	}
	if !strings.Contains(recorder.Body.String(), "MCP_REQUEST_BODY_TOO_LARGE") {
		t.Fatalf("body = %q", recorder.Body.String())
	}
}

func TestLimitMCPRequestBodyAllowsExactLimit(t *testing.T) {
	handler := limitMCPRequestBody(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read exact-limit body: %v", err)
		}
		if string(body) != "12345678" {
			t.Fatalf("body = %q", body)
		}
		w.WriteHeader(http.StatusNoContent)
	}), 8)

	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader("12345678"))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
}

func TestLimitMCPRequestBodyCapsUnknownLengthBody(t *testing.T) {
	var readErr error
	handler := limitMCPRequestBody(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, readErr = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusNoContent)
	}), 8)

	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader("123456789"))
	req.ContentLength = -1
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	var maxErr *http.MaxBytesError
	if !errors.As(readErr, &maxErr) {
		t.Fatalf("read error = %v, want *http.MaxBytesError", readErr)
	}
	if maxErr.Limit != 8 {
		t.Fatalf("limit = %d, want 8", maxErr.Limit)
	}
}

func TestLimitMCPRequestBodyDoesNotLimitNonPost(t *testing.T) {
	handler := limitMCPRequestBody(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read non-POST body: %v", err)
		}
		if string(body) != "123456789" {
			t.Fatalf("body = %q", body)
		}
		w.WriteHeader(http.StatusNoContent)
	}), 8)

	req := httptest.NewRequest(http.MethodDelete, "/mcp", strings.NewReader("123456789"))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
}
