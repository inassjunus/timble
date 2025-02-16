package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	middlewarev1 "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5/middleware"
)

type CtxKey string

const (
	CtxRequestBodyKey = CtxKey("req_body")
	CtxUserIDKey      = CtxKey("user_id")
)

func ReqBodyCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Body != nil {
			buf, _ := io.ReadAll(r.Body)
			defer r.Body.Close()
			ctx = context.WithValue(ctx, CtxRequestBodyKey, string(buf))
			r.Body = io.NopCloser(bytes.NewBuffer(buf))
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ReqIDCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := middleware.GetReqID(r.Context())
		// assign request ID for chizap logger. This line is necessary because of this bug https://github.com/moul/chizap/issues/71
		ctx := context.WithValue(r.Context(), middlewarev1.RequestIDKey, reqId)
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Authentication checks if the user has a valid JWT token
func Authentication(auth *AuthConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				authFailed(w)
				return
			}

			// The token should be prefixed with "Bearer "
			tokenParts := strings.Split(tokenString, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				authFailed(w)
				return
			}

			tokenString = tokenParts[1]

			claims, err := auth.VerifyToken(tokenString)
			if err != nil {
				authFailed(w)
				return
			}

			ctx := context.WithValue(r.Context(), CtxUserIDKey, claims["user_id"])
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func authFailed(w http.ResponseWriter) {
	errByte, _ := json.Marshal(ErrorUnauthenticated)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(errByte)
}
