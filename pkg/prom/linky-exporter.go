package prom

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/linky-exporter/pkg/core"
)

// LinkyExporter object to run exporter server and expose metrics
type LinkyExporter struct {
	Address string
	Port    int
}

// Run method to run http exporter server
func (exporter *LinkyExporter) Run(connector core.LinkyConnector) {
	log.Info(fmt.Sprintf("Beginning to serve on port :%d", exporter.Port))

	prometheus.MustRegister(NewLinkyCollector(connector))
	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", exporter.Address, exporter.Port), nil))
}
