package prometheus

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Exporter struct{}

func New(path string, port int) Exporter {
	http.Handle(path, promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	return Exporter{}
}
