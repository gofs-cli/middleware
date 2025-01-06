package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Expose(r *http.ServeMux) *http.ServeMux {
	r.Handle("GET /metrics", promhttp.Handler())
	return r
}
