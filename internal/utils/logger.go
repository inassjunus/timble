package utils

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func BuildRequestLogFields(r *http.Request, httpStatus int) []zapcore.Field {
	reqBody, _ := r.Context().Value("req_body").(string)
	return []zapcore.Field{
		zap.String("request_id", r.Context().Value(middleware.RequestIDKey).(string)),
		zap.String("req_method", r.Method),
		zap.String("req_url_pattern", chi.RouteContext(r.Context()).RoutePattern()),
		zap.String("req_url", r.URL.String()),
		zap.String("req_body", reqBody),
		zap.Int("resp_http_status", httpStatus),
	}
}
