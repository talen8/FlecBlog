package mcp

import (
	"net/http"
	"sync"
)

const (
	defaultMaxPublicSessions = 128
	mcpSessionIDHeader       = "Mcp-Session-Id"
)

type sessionAdmissionHandler struct {
	next        http.Handler
	maxSessions int

	mu           sync.Mutex
	sessions     map[string]struct{}
	reservations int
}

func newSessionAdmissionHandler(next http.Handler, maxSessions int) http.Handler {
	if maxSessions <= 0 {
		maxSessions = defaultMaxPublicSessions
	}
	return &sessionAdmissionHandler{
		next:        next,
		maxSessions: maxSessions,
		sessions:    make(map[string]struct{}),
	}
}

func (h *sessionAdmissionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get(mcpSessionIDHeader)
	if r.Method == http.MethodPost && sessionID == "" {
		if !h.reserveSessionSlot() {
			http.Error(w, "MCP_SESSION_LIMIT_REACHED", http.StatusServiceUnavailable)
			return
		}
		defer func() {
			h.completeSessionReservation(w.Header().Get(mcpSessionIDHeader))
		}()
		h.next.ServeHTTP(w, r)
		return
	}

	if r.Method == http.MethodDelete && sessionID != "" {
		capture := &statusCaptureResponseWriter{ResponseWriter: w}
		h.next.ServeHTTP(capture, r)
		if capture.status == http.StatusNotFound || capture.status >= 200 && capture.status < 300 {
			h.forgetSession(sessionID)
		}
		return
	}

	h.next.ServeHTTP(w, r)
}

func (h *sessionAdmissionHandler) reserveSessionSlot() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.sessions)+h.reservations >= h.maxSessions {
		return false
	}
	h.reservations++
	return true
}

func (h *sessionAdmissionHandler) completeSessionReservation(sessionID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.reservations > 0 {
		h.reservations--
	}
	if sessionID != "" {
		h.sessions[sessionID] = struct{}{}
	}
}

func (h *sessionAdmissionHandler) forgetSession(sessionID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.sessions, sessionID)
}

type statusCaptureResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusCaptureResponseWriter) WriteHeader(status int) {
	if w.status == 0 {
		w.status = status
	}
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusCaptureResponseWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.ResponseWriter.Write(data)
}
