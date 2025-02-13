package utils

import (
	"encoding/json"
	"net/http"
)

type APIResponser interface {
	ToBytes() []byte
}

type Body struct {
	Meta Meta `json:"meta"`
	Data any  `json:"data"`
}

type MessageBody struct {
	Meta    Meta   `json:"meta"`
	Message string `json:"message"`
}

type Meta struct {
	HTTPStatus int `json:"http_status"`
}

func NewResponseBody(data any, meta Meta) *Body {
	return &Body{
		Data: data,
		Meta: meta,
	}
}

func (b *Body) WriteAPIResponse(w http.ResponseWriter, r *http.Request, httpStatus int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(b.ToBytes())
}

func (b *Body) ToBytes() []byte {
	res, _ := json.Marshal(b)

	return res
}

func NewResponseMessage(message string, meta Meta) *MessageBody {
	return &MessageBody{
		Meta:    meta,
		Message: message,
	}
}

func (b *MessageBody) WriteAPIResponse(w http.ResponseWriter, r *http.Request, httpStatus int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(b.ToBytes())
}

func (b *MessageBody) ToBytes() []byte {
	res, _ := json.Marshal(b)

	return res
}
