package core

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/linky-exporter/pkg/collectors"
	"go.bug.st/serial"
)

// LinkyExporter object to run exporter server and expose metrics
type LinkyExporter struct {
	Address   string
	Port      int
	Device    string
	BaudRate  int
	FrameSize int
	Parity    string
	StopBits  string
}

// Run method to run http exporter server
func (exporter *LinkyExporter) Run() {
	log.Info(fmt.Sprintf("Beginning to serve on port :%d", exporter.Port))

	prometheus.MustRegister(collectors.NewLinkyCollector(exporter.Device, exporter.BaudRate,
		exporter.FrameSize, parseParity(exporter.Parity), parseStopBits(exporter.StopBits)))
	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", exporter.Address, exporter.Port), nil))
}

func parseParity(value string) (parity serial.Parity) {
	switch value {
	case "ParityNone", "N":
		parity = serial.NoParity
		break
	case "ParityOdd", "O":
		parity = serial.OddParity
		break
	case "ParityEven", "E":
		parity = serial.EvenParity
		break
	case "ParityMark", "M":
		parity = serial.MarkParity
		break
	case "ParitySpace", "S":
		parity = serial.SpaceParity
		break
	default:
		_, err := fmt.Fprintln(os.Stderr, "Impossible to parse Parity named", value)
		if err != nil {
			log.Error(err)
		}
		os.Exit(3)
	}
	return
}

func parseStopBits(value string) (stopBits serial.StopBits) {
	switch value {
	case "Stop1", "1":
		stopBits = serial.OneStopBit
		break
	case "Stop1Half", "15":
		stopBits = serial.OnePointFiveStopBits
		break
	case "Stop2", "2":
		stopBits = serial.TwoStopBits
		break
	default:
		_, err := fmt.Fprintln(os.Stderr, "Impossible to parse StopBits named", value)
		if err != nil {
			log.Error(err)
		}
		os.Exit(3)
	}
	return
}
