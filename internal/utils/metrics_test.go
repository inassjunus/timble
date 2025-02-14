package utils_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"

	"timble/internal/utils"
)

const (
	testClientName = "test_service"
	testActionName = "test_action"
)

func TestRegisterCustomMetrics(t *testing.T) {
	t.Run("register metrics", func(t *testing.T) {
		assert.NotPanics(t, func() {
			utils.RegisterCustomMetrics(prometheus.NewRegistry())
		})
	})
}

func TestClientMetric_NewClientMetric(t *testing.T) {
	cases := []struct {
		name       string
		clientName string
		actionName string
	}{
		{
			name:       "normal case",
			clientName: testClientName,
			actionName: testActionName,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := utils.NewClientMetric(tc.clientName, tc.actionName)

			assert.Equal(t, testClientName, actual.ClientName)
			assert.Equal(t, testActionName, actual.Action)
			assert.Equal(t, utils.RequestStatusOK, actual.RequestStatus)
			assert.Equal(t, "", actual.HTTPStatus)
			assert.NotNil(t, actual.StartTime)
		})
	}
}

func TestClientMetric_TrackClient(t *testing.T) {
	cases := []struct {
		name   string
		metric *utils.ClientMetric
	}{
		{
			name:   "track basic metric case",
			metric: utils.NewClientMetric(testClientName, testActionName),
		},
		{
			name:   "track successful request case",
			metric: utils.NewClientMetric(testClientName, testActionName).SetHttpStatus(200),
		},
		{
			name:   "track failed request case",
			metric: utils.NewClientMetric(testClientName, testActionName).SetFail(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				tc.metric.TrackClient()
			})
		})
	}
}

func TestClientMetric_TrackClientWithError(t *testing.T) {
	cases := []struct {
		name   string
		metric *utils.ClientMetric
		err    error
	}{
		{
			name:   "track non-error metric case",
			metric: utils.NewClientMetric(testClientName, testActionName),
		},
		{
			name:   "track error metric case",
			metric: utils.NewClientMetric(testClientName, testActionName),
			err:    errors.New("timeout"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				tc.metric.TrackClientWithError(tc.err)
			})
		})
	}
}

func TestClientMetric_SetFail(t *testing.T) {
	cases := []struct {
		name       string
		clientName string
		actionName string
	}{
		{
			name:       "normal case",
			clientName: testClientName,
			actionName: testActionName,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := utils.NewClientMetric(tc.clientName, tc.actionName).SetFail()

			assert.Equal(t, utils.RequestStatusFail, actual.RequestStatus)
			// also ensure that other values remain the same
			assert.Equal(t, testClientName, actual.ClientName)
			assert.Equal(t, testActionName, actual.Action)
			assert.Equal(t, "", actual.HTTPStatus)
			assert.NotNil(t, actual.StartTime)
		})
	}
}

func TestClientMetric_SetHttpStatus(t *testing.T) {
	cases := []struct {
		name       string
		clientName string
		actionName string
	}{
		{
			name:       "normal case",
			clientName: testClientName,
			actionName: testActionName,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := utils.NewClientMetric(tc.clientName, tc.actionName).SetHttpStatus(200)

			assert.Equal(t, "200", actual.HTTPStatus)
			// also ensure that other values remain the same
			assert.Equal(t, testClientName, actual.ClientName)
			assert.Equal(t, testActionName, actual.Action)
			assert.Equal(t, utils.RequestStatusOK, actual.RequestStatus)
			assert.NotNil(t, actual.StartTime)
		})
	}
}

func TestRestMetric_NewRestMetric(t *testing.T) {
	cases := []struct {
		name   string
		method string
	}{
		{
			name:   "normal case",
			method: http.MethodGet,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), chi.RouteCtxKey, chi.NewRouteContext())
			req := httptest.NewRequest(tc.method, "/", bytes.NewBuffer(nil)).WithContext(ctx)
			actual := utils.NewRestMetric(req)

			assert.Equal(t, tc.method, actual.RequestMethod)
			assert.Equal(t, "", actual.RequestUrl)
			assert.Equal(t, utils.RequestStatusOK, actual.RequestStatus)
			assert.Equal(t, http.StatusOK, actual.HTTPStatus)
			assert.NotNil(t, actual.StartTime)
		})
	}
}

func TestRestMetric_SetFail(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		httpStatus int
	}{
		{
			name:       "normal case",
			method:     http.MethodGet,
			httpStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), chi.RouteCtxKey, chi.NewRouteContext())
			req := httptest.NewRequest(tc.method, "/", bytes.NewBuffer(nil)).WithContext(ctx)
			actual := utils.NewRestMetric(req).SetFail(tc.httpStatus)

			assert.Equal(t, utils.RequestStatusFail, actual.RequestStatus)
			assert.Equal(t, tc.httpStatus, actual.HTTPStatus)

			// also ensure that other values remain the same
			assert.Equal(t, tc.method, actual.RequestMethod)
			assert.Equal(t, "", actual.RequestUrl)
			assert.NotNil(t, actual.StartTime)
		})
	}
}

func TestMetrics_TrackRestService(t *testing.T) {
	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, chi.NewRouteContext())
	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(nil)).WithContext(ctx)
	cases := []struct {
		name   string
		metric *utils.RestMetric
	}{
		{
			name:   "track successful request case",
			metric: utils.NewRestMetric(req),
		},
		{
			name:   "track failed request case",
			metric: utils.NewRestMetric(req).SetFail(http.StatusUnauthorized),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				tc.metric.TrackRestService()
			})
		})
	}
}
