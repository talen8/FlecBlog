package mcp

import "net/http"

const defaultMaxMCPRequestBodyBytes int64 = 8 << 20

func limitMCPRequestBody(next http.Handler, maxBytes int64) http.Handler {
	if maxBytes <= 0 {
		maxBytes = defaultMaxMCPRequestBodyBytes
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.Body == nil {
			next.ServeHTTP(w, r)
			return
		}
		if r.ContentLength > maxBytes {
			http.Error(w, "MCP_REQUEST_BODY_TOO_LARGE", http.StatusRequestEntityTooLarge)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		next.ServeHTTP(w, r)
	})
}
