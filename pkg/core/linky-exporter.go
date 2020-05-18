package core

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/syberalexis/linky-exporter/pkg/collectors"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type LinkyExporter struct {
	Port int
}

func (l *LinkyExporter) Run() {
	prometheus.MustRegister(collectors.NewLinkyCollector())

	http.Handle("/metrics", promhttp.Handler())
	log.Info(fmt.Sprintf("Beginning to serve on port :%d", l.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", l.Port), nil))
}
