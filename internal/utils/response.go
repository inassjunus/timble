package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Meta        Meta           `json:"meta"`
	Data        any            `json:"data,omitempty"`
	Error       string         `json:"error,omitempty"`
	ErrorDetail *StandardError `json:"error_detail,omitempty"`
	Message     string         `json:"message,omitempty"`
}

type Meta struct {
	HTTPStatus int `json:"http_status"`
}

func NewDataResponse(data any, meta Meta) *Response {
	return &Response{
		Data: data,
		Meta: meta,
	}
}

func NewMessageResponse(message string, meta Meta) *Response {
	return &Response{
		Meta:    meta,
		Message: message,
	}
}

func NewErrorResponse(err error, httpStatus int) *Response {
	result := &Response{
		Error: err.Error(),
		Meta: Meta{
			HTTPStatus: httpStatus,
		},
	}
	errStd, ok := err.(*StandardError)
	if ok {
		result.Error = ""
		result.ErrorDetail = errStd
	}
	return result
}

func (b *Response) WriteAPIResponse(w http.ResponseWriter, r *http.Request, httpStatus int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(b.ToBytes())
}

func (b *Response) ToBytes() []byte {
	res, _ := json.Marshal(b)

	return res
}
