package utils_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"timble/internal/utils"
)

func TestLogger_BuildRequestLogFields(t *testing.T) {
	testBody := `{
    "tes": [1, 2, 3]
  }`
	type req struct {
		method string
		url    string
		body   io.Reader
		route  string
	}
	cases := []struct {
		name       string
		req        req
		httpStatus int
		expected   []zapcore.Field
	}{
		{
			name: "normal case with blank body",
			req: req{
				method: http.MethodGet,
				url:    "/test/123?p=1",
				route:  "/test/{id}",
			},
			httpStatus: http.StatusOK,
			expected: []zapcore.Field{
				zap.String("request_id", ""),
				zap.String("req_method", http.MethodGet),
				zap.String("req_url_pattern", "/test/{id}"),
				zap.String("req_url", "/test/123?p=1"),
				zap.String("req_body", ""),
				zap.Int("resp_http_status", http.StatusOK),
			},
		},
		{
			name: "normal case with body",
			req: req{
				method: http.MethodPost,
				url:    "/test/123?p=1",
				body:   io.NopCloser(strings.NewReader(testBody)),
				route:  "/test/{id}",
			},
			httpStatus: http.StatusOK,
			expected: []zapcore.Field{
				zap.String("request_id", ""),
				zap.String("req_method", http.MethodPost),
				zap.String("req_url_pattern", "/test/{id}"),
				zap.String("req_url", "/test/123?p=1"),
				zap.String("req_body", "{\n    \"tes\": [1, 2, 3]\n  }"),
				zap.Int("resp_http_status", http.StatusOK),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			chiCtx := chi.NewRouteContext()
			chiCtx.RoutePatterns = []string{c.req.route}
			ctx := context.WithValue(context.Background(), chi.RouteCtxKey, chiCtx)
			req, _ := http.NewRequest(c.req.method, c.req.url, c.req.body)

			if req.Body != nil {
				buf, _ := io.ReadAll(req.Body)
				defer req.Body.Close()
				ctx = context.WithValue(ctx, utils.CtxRequestBodyKey, string(buf))
				req.Body = io.NopCloser(bytes.NewBuffer(buf))
			}

			actual := utils.BuildRequestLogFields(req.WithContext(ctx), c.httpStatus)

			assert.Equal(t, c.expected, actual)
		})
	}
}
