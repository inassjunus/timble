package utils

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/prometheus/client_golang/prometheus"
)

var TimbleRestServiceDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "timble_rest_service_duration",
	Help: "track request duration of search tuning request",
}, []string{"method", "url", "status", "http_status"})

var TimbleClientDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "timble_duration_for_requesting_to_client",
	Help: "track request duration to client",
}, []string{"client_name", "action", "status", "http_status"})

const (
	RequestStatusOK   = "ok"
	RequestStatusFail = "fail"
)

type ClientMetric struct {
	ClientName    string
	Action        string
	StartTime     time.Time
	RequestStatus string
	HTTPStatus    string
}

type RestMetric struct {
	RequestMethod string
	RequestUrl    string
	StartTime     time.Time
	RequestStatus string
	HTTPStatus    int
}

func RegisterCustomMetrics(reg *prometheus.Registry) {
	reg.MustRegister(
		TimbleRestServiceDuration,
		TimbleClientDuration,
	)
}

func NewClientMetric(clientName, action string) *ClientMetric {
	return &ClientMetric{
		ClientName:    clientName,
		Action:        action,
		StartTime:     time.Now(),
		RequestStatus: RequestStatusOK,
	}
}

func (m *ClientMetric) TrackClient() {
	elapsed := time.Since(m.StartTime)
	TimbleClientDuration.WithLabelValues(m.ClientName, m.Action, m.RequestStatus, m.HTTPStatus).Observe(elapsed.Seconds())
}

func (m *ClientMetric) SetFail() *ClientMetric {
	m.RequestStatus = RequestStatusFail
	return m
}

func (m *ClientMetric) SetHttpStatus(status int) *ClientMetric {
	m.HTTPStatus = strconv.Itoa(status)
	return m
}

func (m *ClientMetric) TrackClientWithError(err error) {
	if err != nil {
		m.RequestStatus = RequestStatusFail
	}
	m.TrackClient()
}

func NewRestMetric(r *http.Request) *RestMetric {
	return &RestMetric{
		RequestMethod: r.Method,
		RequestUrl:    chi.RouteContext(r.Context()).RoutePattern(),
		StartTime:     time.Now(),
		RequestStatus: RequestStatusOK,
		HTTPStatus:    http.StatusOK,
	}
}

func (m *RestMetric) SetFail(httpStatus int) *RestMetric {
	m.RequestStatus = RequestStatusFail
	m.HTTPStatus = httpStatus
	return m
}

func (m *RestMetric) TrackRestService() {
	elapsed := time.Since(m.StartTime)
	TimbleRestServiceDuration.WithLabelValues(m.RequestMethod, m.RequestUrl, m.RequestStatus, strconv.Itoa(m.HTTPStatus)).Observe(elapsed.Seconds())
}
