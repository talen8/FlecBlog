package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"flec_blog/config"

	"github.com/gin-gonic/gin"
)

func TestCORSNavigationPathAllowsExactGETWithoutGrantingCrossOriginRead(t *testing.T) {
	gin.SetMode(gin.TestMode)
	conf := &config.Config{Server: config.ServerConfig{AllowOrigins: []string{"https://aa.example.test"}}}
	router := gin.New()
	router.Use(CORSWithNavigationPaths(conf, "/authorize"))
	router.GET("/authorize", func(c *gin.Context) {
		c.String(http.StatusOK, "reached")
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/authorize?client_id=test", nil)
	request.Header.Set("Origin", "https://chat.example.test")
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("GET /authorize status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if recorder.Body.String() != "reached" {
		t.Fatalf("GET /authorize body = %q", recorder.Body.String())
	}
	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("GET /authorize Access-Control-Allow-Origin = %q, want empty", got)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Credentials"); got != "" {
		t.Fatalf("GET /authorize Access-Control-Allow-Credentials = %q, want empty", got)
	}
}

func TestCORSNavigationPathAllowsSameOriginDocumentPOSTWithoutGrantingCrossOriginRead(t *testing.T) {
	gin.SetMode(gin.TestMode)
	conf := &config.Config{Server: config.ServerConfig{AllowOrigins: []string{"https://aa.example.test"}}}
	router := gin.New()
	router.Use(CORSWithNavigationPaths(conf, "/authorize"))
	router.POST("/authorize", func(c *gin.Context) {
		c.String(http.StatusOK, "reached")
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/authorize", nil)
	request.Header.Set("Origin", "null")
	request.Header.Set("Sec-Fetch-Mode", "navigate")
	request.Header.Set("Sec-Fetch-Dest", "document")
	request.Header.Set("Sec-Fetch-Site", "same-origin")
	request.Header.Set("Sec-Fetch-User", "?1")
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("POST /authorize status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if recorder.Body.String() != "reached" {
		t.Fatalf("POST /authorize body = %q", recorder.Body.String())
	}
	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("POST /authorize Access-Control-Allow-Origin = %q, want empty", got)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Credentials"); got != "" {
		t.Fatalf("POST /authorize Access-Control-Allow-Credentials = %q, want empty", got)
	}
}

func TestCORSNavigationPathStillRejectsUnallowlistedOriginForOtherRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	conf := &config.Config{Server: config.ServerConfig{AllowOrigins: []string{"https://aa.example.test"}}}
	router := gin.New()
	router.Use(CORSWithNavigationPaths(conf, "/authorize"))
	router.GET("/api/v1/articles", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.POST("/authorize", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.GET("/authorize/extra", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	tests := []struct {
		name    string
		method  string
		path    string
		headers map[string]string
	}{
		{name: "unrelated API", method: http.MethodGet, path: "/api/v1/articles"},
		{name: "authorization POST without navigation metadata", method: http.MethodPost, path: "/authorize"},
		{
			name:   "cross-site authorization POST navigation",
			method: http.MethodPost,
			path:   "/authorize",
			headers: map[string]string{
				"Sec-Fetch-Mode": "navigate",
				"Sec-Fetch-Dest": "document",
				"Sec-Fetch-Site": "cross-site",
			},
		},
		{
			name:   "same-origin authorization POST non-document",
			method: http.MethodPost,
			path:   "/authorize",
			headers: map[string]string{
				"Sec-Fetch-Mode": "cors",
				"Sec-Fetch-Dest": "empty",
				"Sec-Fetch-Site": "same-origin",
			},
		},
		{name: "non-exact authorization path", method: http.MethodGet, path: "/authorize/extra"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(tt.method, tt.path, nil)
			request.Header.Set("Origin", "https://chat.example.test")
			for name, value := range tt.headers {
				request.Header.Set(name, value)
			}
			router.ServeHTTP(recorder, request)
			if recorder.Code != http.StatusForbidden {
				t.Fatalf("%s %s status = %d, want %d", tt.method, tt.path, recorder.Code, http.StatusForbidden)
			}
		})
	}
}

func TestCORSNavigationPathPreservesAllowlistedCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const allowedOrigin = "https://aa.example.test"
	conf := &config.Config{Server: config.ServerConfig{AllowOrigins: []string{allowedOrigin}}}
	router := gin.New()
	router.Use(CORSWithNavigationPaths(conf, "/authorize"))
	router.GET("/authorize", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/authorize", nil)
	request.Header.Set("Origin", allowedOrigin)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("allowlisted GET /authorize status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != allowedOrigin {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, allowedOrigin)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Fatalf("Access-Control-Allow-Credentials = %q, want true", got)
	}
}
