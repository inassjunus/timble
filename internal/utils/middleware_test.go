package utils_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"

	"timble/internal/utils"
)

func TestMiddleware_ReqIDCtx(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		reqId := middleware.GetReqID(r.Context())
		contentType := w.Header().Get("Content-Type")
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("reqID: %s, contentType: %s", reqId, contentType)))
	}

	router := chi.NewRouter()
	router.Route("/test", func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(utils.ReqIDCtx)
		r.Get("/", testHandler)
	})
	ts := httptest.NewServer(router)

	cases := []struct {
		name     string
		reqID    string
		expected string
	}{
		{
			name:     "normal case request ID is given on header",
			reqID:    "test/abc0001",
			expected: `reqID: test\/abc0001, contentType: application\/json`,
		},
		{
			name: "normal case - request ID is not given on header",
			// request id format from chi/v5/middleware is "<hostname>/<random 10 chars>-<6 digit number>
			expected: `reqID: (.+\/.{10}-\d{6}), contentType: application\/json`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, body := testRequestWithReqID(t, ts, "GET", "/test", nil, c.reqID)
			match, _ := regexp.MatchString(c.expected, body)
			assert.Equal(t, true, match)
		})
	}
	defer ts.Close()
}

func TestMiddleware_ReqBodyCtx(t *testing.T) {
	testBody := `{
    "user_ids": [1, 2, 3],
    "enabled": true
  }`
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		reqBodyFromCtx := "N/A"
		reqBody := "N/A"
		if r.Body != nil {
			buf, _ := io.ReadAll(r.Body)
			defer r.Body.Close()
			reqBody = string(buf)
			reqBodyFromCtx, _ = r.Context().Value(utils.CtxRequestBodyKey).(string)
		}
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("req body from ctx: %s, req body from request: %s", reqBodyFromCtx, reqBody)))
	}

	router := chi.NewRouter()
	router.Route("/test", func(r chi.Router) {
		r.Use(utils.ReqBodyCtx)
		r.Post("/post", testHandler)
		r.Get("/get", testHandler)
	})
	ts := httptest.NewServer(router)

	cases := []struct {
		name     string
		reqBody  io.Reader
		url      string
		method   string
		expected string
	}{
		{
			name:     "normal case when request body is given",
			url:      "/test/post",
			method:   http.MethodPost,
			reqBody:  io.NopCloser(strings.NewReader(testBody)),
			expected: "req body from ctx: {\n    \"user_ids\": [1, 2, 3],\n    \"enabled\": true\n  }, req body from request: {\n    \"user_ids\": [1, 2, 3],\n    \"enabled\": true\n  }",
		},
		{
			name:     "normal case when request body is not given",
			url:      "/test/get",
			method:   http.MethodGet,
			expected: "req body from ctx: , req body from request: ",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, body := testRequest(t, ts, c.method, c.url, c.reqBody)
			assert.Equal(t, c.expected, body)
		})
	}
	defer ts.Close()
}

func TestMiddleware_Authentication(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}

	cfg := &utils.AuthConfig{
		SecretKey: []byte("secretz"),
		TokenExp:  time.Hour,
	}

	router := chi.NewRouter()
	router.Route("/test", func(r chi.Router) {
		r.Use(utils.Authentication(cfg))
		r.Get("/", testHandler)
	})
	ts := httptest.NewServer(router)

	cases := []struct {
		name               string
		token              string
		expectedResult     string
		expectedHTTPStatus int
	}{
		{
			name:               "normal case",
			expectedResult:     "OK",
			expectedHTTPStatus: 200,
		},
		{
			name:               "missing token case",
			expectedResult:     `{"message":"Invalid or missing required authentication","code":"Unauthorized"}`,
			expectedHTTPStatus: 401,
		},
		{
			name:               "invalid token format",
			token:              "this is invalid header",
			expectedResult:     `{"message":"Invalid or missing required authentication","code":"Unauthorized"}`,
			expectedHTTPStatus: 401,
		},
		{
			name:               "invalid token case",
			token:              "thisisinvalidtoken",
			expectedResult:     `{"message":"Invalid or missing required authentication","code":"Unauthorized"}`,
			expectedHTTPStatus: 401,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/test", ts.URL), nil)
			if err != nil {
				t.Fatal(err)
			}
			if tc.expectedHTTPStatus == http.StatusOK {
				tc.token, _ = cfg.GenerateToken(uint(1))
			}
			if tc.token != "" {
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tc.token))
			}
			res, body := testCallRequest(t, ts, req)
			assert.Equal(t, tc.expectedResult, body)
			assert.Equal(t, tc.expectedHTTPStatus, res.StatusCode)
		})
	}
	defer ts.Close()
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	return testRequestWithReqID(t, ts, method, path, body, "test/0000001")
}

func testRequestWithReqID(t *testing.T, ts *httptest.Server, method, path string, body io.Reader, reqID string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	req.Header.Set("X-Request-Id", reqID)
	return testCallRequest(t, ts, req)
}

func testCallRequest(t *testing.T, ts *httptest.Server, req *http.Request) (*http.Response, string) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}
