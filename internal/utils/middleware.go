package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	middlewarev1 "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
)

func ReqBodyCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Body != nil {
			buf, _ := io.ReadAll(r.Body)
			defer r.Body.Close()
			ctx = context.WithValue(ctx, "req_body", string(buf))
			r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
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
				authFailed(w, errors.New("Missing token"))
				return
			}

			// The token should be prefixed with "Bearer "
			tokenParts := strings.Split(tokenString, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				authFailed(w, errors.New("Invalid token"))
				return
			}

			tokenString = tokenParts[1]

			claims, err := auth.VerifyToken(tokenString)
			if err != nil {
				authFailed(w, errors.New("Invalid token"))
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func authFailed(w http.ResponseWriter, err error) {
	errByte, _ := json.Marshal(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(errByte)
}
