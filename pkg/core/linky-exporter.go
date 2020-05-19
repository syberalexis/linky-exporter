package core

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/linky-exporter/pkg/collectors"
	"net/http"
)

type LinkyExporter struct {
	Address  string
	Port     int
	File     string
	BaudRate int
}

func (exporter *LinkyExporter) Run() {
	prometheus.MustRegister(collectors.NewLinkyCollector(exporter.File, exporter.BaudRate))

	log.Info(fmt.Sprintf("Beginning to serve on port :%d", exporter.Port))
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", exporter.Address, exporter.Port), nil))
}
