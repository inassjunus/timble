package utils_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"timble/internal/utils"
)

func TestResponse_NewDataResponse(t *testing.T) {
	cases := []struct {
		name     string
		body     any
		meta     utils.Meta
		expected *utils.Response
	}{
		{
			name: "normal case",
			body: "Everything's fine",
			meta: utils.Meta{HTTPStatus: 200},
			expected: &utils.Response{
				Data: "Everything's fine",
				Meta: utils.Meta{HTTPStatus: 200},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := utils.NewDataResponse(c.body, c.meta)

			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestResponse_NewMessageResponse(t *testing.T) {
	cases := []struct {
		name     string
		message  string
		meta     utils.Meta
		expected *utils.Response
	}{
		{
			name:    "normal case",
			message: "Everything's fine",
			meta:    utils.Meta{HTTPStatus: 200},
			expected: &utils.Response{
				Message: "Everything's fine",
				Meta:    utils.Meta{HTTPStatus: 200},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := utils.NewMessageResponse(c.message, c.meta)

			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestResponse_NewErrorResponse(t *testing.T) {
	cases := []struct {
		name       string
		err        error
		httpStatus int
		expected   *utils.Response
	}{
		{
			name:       "normal case with default error",
			err:        errors.New("Unexpected error"),
			httpStatus: http.StatusUnauthorized,
			expected: &utils.Response{
				Error: "Unexpected error",
				Meta:  utils.Meta{HTTPStatus: http.StatusUnauthorized},
			},
		},
		{
			name:       "normal case with standard error",
			err:        utils.ErrUnauthenticated,
			httpStatus: http.StatusUnauthorized,
			expected: &utils.Response{
				ErrorDetail: utils.ErrUnauthenticated,
				Meta:        utils.Meta{HTTPStatus: http.StatusUnauthorized},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := utils.NewErrorResponse(tc.err, http.StatusUnauthorized)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestResponse_WriteAPIResponse(t *testing.T) {
	cases := []struct {
		name     string
		response utils.Response
		meta     utils.Meta
	}{
		{
			name: "normal case with data",
			response: utils.Response{
				Data: []string{"Everything's fine"},
				Meta: utils.Meta{HTTPStatus: 200},
			},
			meta: utils.Meta{HTTPStatus: 200},
		},
		{
			name: "normal case with message",
			response: utils.Response{
				Message: "Everything's fine",
				Meta:    utils.Meta{HTTPStatus: 200},
			},
			meta: utils.Meta{HTTPStatus: 200},
		},
		{
			name: "normal case with error",
			response: utils.Response{
				Error: "Everything's NOT fine",
				Meta:  utils.Meta{HTTPStatus: 500},
			},
			meta: utils.Meta{HTTPStatus: 500},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				tc.response.WriteAPIResponse(httptest.NewRecorder(), &http.Request{}, 200)
			})
		})
	}
}

func TestBodyResponse_ToBytes(t *testing.T) {
	cases := []struct {
		name     string
		response utils.Response
		expected []byte
	}{
		{
			name: "normal case with data",
			response: utils.Response{
				Data: []string{"Everything's fine"},
				Meta: utils.Meta{HTTPStatus: 200},
			},
			expected: []byte("{\"meta\":{\"http_status\":200},\"data\":[\"Everything's fine\"]}"),
		},
		{
			name: "normal case with message",
			response: utils.Response{
				Message: "Everything's fine",
				Meta:    utils.Meta{HTTPStatus: 200},
			},
			expected: []byte("{\"meta\":{\"http_status\":200},\"message\":\"Everything's fine\"}"),
		},
		{
			name: "normal case with error",
			response: utils.Response{
				Error: "Everything's NOT fine",
				Meta:  utils.Meta{HTTPStatus: 400},
			},
			expected: []byte("{\"meta\":{\"http_status\":400},\"error\":\"Everything's NOT fine\"}"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.response.ToBytes()
			assert.Equal(t, tc.expected, actual)
		})
	}
}
