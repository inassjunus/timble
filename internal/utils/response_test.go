package utils_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"timble/internal/utils"
)

func TestResponse_NewResponseBody(t *testing.T) {
	cases := []struct {
		name     string
		body     any
		meta     utils.Meta
		expected *utils.Body
	}{
		{
			name: "normal case",
			body: "Everything's fine",
			meta: utils.Meta{HTTPStatus: 200},
			expected: &utils.Body{
				Data: "Everything's fine",
				Meta: utils.Meta{HTTPStatus: 200},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := utils.NewResponseBody(c.body, c.meta)

			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestResponseBody_WriteAPIResponse(t *testing.T) {
	cases := []struct {
		name string
		body any
		meta utils.Meta
	}{
		{
			name: "normal case",
			body: "Everything's fine",
			meta: utils.Meta{HTTPStatus: 200},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			writer := utils.NewResponseBody(c.body, c.meta)
			assert.NotPanics(t, func() {
				writer.WriteAPIResponse(httptest.NewRecorder(), &http.Request{}, 200)
			})
		})
	}
}

func TestResponseBody_ToBytes(t *testing.T) {
	cases := []struct {
		name     string
		body     any
		meta     utils.Meta
		expected []byte
	}{
		{
			name:     "normal case",
			body:     "Everything's fine",
			meta:     utils.Meta{HTTPStatus: 200},
			expected: []byte("{\"meta\":{\"http_status\":200},\"data\":\"Everything's fine\"}"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := utils.NewResponseBody(c.body, c.meta).ToBytes()
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestResponse_NewResponseMessage(t *testing.T) {
	cases := []struct {
		name     string
		message  string
		meta     utils.Meta
		expected *utils.MessageBody
	}{
		{
			name:    "normal case",
			message: "Everything's fine",
			meta:    utils.Meta{HTTPStatus: 200},
			expected: &utils.MessageBody{
				Message: "Everything's fine",
				Meta:    utils.Meta{HTTPStatus: 200},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := utils.NewResponseMessage(c.message, c.meta)

			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestResponseMessageBody_WriteAPIResponse(t *testing.T) {
	cases := []struct {
		name string
		body string
		meta utils.Meta
	}{
		{
			name: "normal case",
			body: "Everything's fine",
			meta: utils.Meta{HTTPStatus: 200},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			writer := utils.NewResponseMessage(c.body, c.meta)
			assert.NotPanics(t, func() {
				writer.WriteAPIResponse(httptest.NewRecorder(), &http.Request{}, 200)
			})
		})
	}
}

func TestResponseMessageBody_ToBytes(t *testing.T) {
	cases := []struct {
		name     string
		body     string
		meta     utils.Meta
		expected []byte
	}{
		{
			name:     "normal case",
			body:     "Everything's fine",
			meta:     utils.Meta{HTTPStatus: 200},
			expected: []byte("{\"meta\":{\"http_status\":200},\"message\":\"Everything's fine\"}"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := utils.NewResponseMessage(c.body, c.meta).ToBytes()
			assert.Equal(t, c.expected, actual)
		})
	}
}
