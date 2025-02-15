package config_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"

	"timble/internal/config"
)

func Test_NewRestServer(t *testing.T) {
	messageResponseBase := `{
    "meta":{
        "http_status":%d
    },
    "message":"%s"
  }`
	rds := miniredis.RunT(t)
	type expected struct {
		expectedResponse   string
		expectedHTTPStatus int
		expectedPanic      bool
		expectedErr        error
	}
	tests := []struct {
		name     string
		mockFn   func()
		expected expected
	}{
		{
			name: "normal case",
			mockFn: func() {
				os.Setenv("TOKEN_EXPIRATION", "10m")
				os.Setenv("REDIS_HOST", rds.Host())
				os.Setenv("REDIS_PORT", rds.Port())
				os.Setenv("REDIS_TIMEOUT", "200ms")
				os.Setenv("REDIS_DB", "0")
				os.Setenv("CACHE_HOST", rds.Host())
				os.Setenv("CACHE_PORT", rds.Port())
				os.Setenv("CACHE_TIMEOUT", "200ms")
				os.Setenv("CACHE_DB", "0")
			},
			expected: expected{
				expectedHTTPStatus: http.StatusOK,
				expectedResponse:   fmt.Sprintf(messageResponseBase, http.StatusOK, "ok"),
			},
		},
		{
			name: "redis is down",
			mockFn: func() {
				os.Setenv("TOKEN_EXPIRATION", "10m")
				os.Setenv("REDIS_HOST", "")
				os.Setenv("REDIS_PORT", "")
				os.Setenv("REDIS_TIMEOUT", "200ms")
				os.Setenv("REDIS_DB_CACHE", "0")
				os.Setenv("REDIS_DB_STORAGE", "1")
			},
			expected: expected{
				expectedPanic: true,
			},
		},
		{
			name: "cache is down",
			mockFn: func() {
				os.Setenv("TOKEN_EXPIRATION", "10m")
				os.Setenv("REDIS_HOST", rds.Host())
				os.Setenv("REDIS_PORT", rds.Port())
				os.Setenv("REDIS_TIMEOUT", "200ms")
				os.Setenv("REDIS_DB", "0")
				os.Setenv("CACHE_HOST", "")
				os.Setenv("CACHE_PORT", "")
				os.Setenv("CACHE_TIMEOUT", "200ms")
				os.Setenv("CACHE_DB", "0")
			},
			expected: expected{
				expectedPanic: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFn()
			if tc.expected.expectedPanic {
				assert.Panics(t, func() {
					config.NewRestServer()
				})
				return
			}

			server, err := config.NewRestServer()

			if tc.expected.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expected.expectedErr.Error(), err.Error())
			} else {
				assert.NotEqual(t, nil, server)
        // test healthcheck endpoint
				urlPath := "/health"

				req := httptest.NewRequest(http.MethodGet, urlPath, bytes.NewBuffer(nil))
				recorder := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext())
				req = req.WithContext(ctx)

				hndlr := server.Server.Handler
				hndlr.ServeHTTP(recorder, req)

				assert.Equal(t, tc.expected.expectedHTTPStatus, recorder.Result().StatusCode)
				assert.JSONEq(t, tc.expected.expectedResponse, recorder.Body.String())
			}
		})
	}
}
