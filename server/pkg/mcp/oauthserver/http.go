package oauthserver

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

type oauthErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	writeJSONStatus(w, status, payload)
}

func writeJSONStatus(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeOAuthError(w http.ResponseWriter, status int, code, description string) {
	writeJSON(w, status, oauthErrorResponse{Error: code, ErrorDescription: description})
}

func readLimitedJSON(r *http.Request, out any) error {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		return errors.New("content type must be application/json")
	}
	reader := io.LimitReader(r.Body, maxRequestBodyBytes+1)
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(out); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return errors.New("request body must contain one JSON value")
	}
	return nil
}

func readLimitedForm(r *http.Request) (url.Values, error) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/x-www-form-urlencoded" {
		return nil, errors.New("content type must be application/x-www-form-urlencoded")
	}
	reader := io.LimitReader(r.Body, maxRequestBodyBytes+1)
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if len(data) > maxRequestBodyBytes {
		return nil, errors.New("request body too large")
	}
	return url.ParseQuery(string(data))
}

func secureHTMLHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; form-action 'self'; base-uri 'none'; frame-ancestors 'none'")
	w.Header().Set("Referrer-Policy", "no-referrer")
}

func secureAuthorizationHTMLHeaders(w http.ResponseWriter, redirectURI string) {
	secureHTMLHeaders(w)
	origin, ok := cspRedirectOrigin(redirectURI)
	if !ok {
		return
	}
	w.Header().Set(
		"Content-Security-Policy",
		"default-src 'none'; style-src 'unsafe-inline'; form-action 'self' "+origin+"; base-uri 'none'; frame-ancestors 'none'",
	)
}

func cspRedirectOrigin(raw string) (string, bool) {
	parsed, err := url.Parse(raw)
	if err != nil || validateRedirectURI(raw) != nil {
		return "", false
	}
	origin := parsed.Scheme + "://" + parsed.Host
	for _, r := range origin {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || strings.ContainsRune("-._~:/[]", r) {
			continue
		}
		return "", false
	}
	return origin, true
}

func normalizeSpace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
