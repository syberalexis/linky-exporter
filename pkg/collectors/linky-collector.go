package collectors

import (
	"bufio"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
	"strconv"
	"strings"
)

// LinkyCollector object to describe and collect metrics
type LinkyCollector struct {
	device   string
	baudRate int
	optarif  *prometheus.Desc
	imax     *prometheus.Desc
	hchc     *prometheus.Desc
	iinst    *prometheus.Desc
	papp     *prometheus.Desc
	motdetat *prometheus.Desc
	hhphc    *prometheus.Desc
	isousc   *prometheus.Desc
	hchp     *prometheus.Desc
	ptec     *prometheus.Desc
}

// Internal linky values object to each metrics
type linkyValues struct {
	optarif  string // 'OPTARIF': 'HC..',        # option tarifaire
	imax     uint16 // 'IMAX': '007',            # intensité max
	hchc     uint64 // 'HCHC': '040177099',      # index heure creuse en Wh
	iinst    int16  // 'IINST': '005',           # Intensité instantanée en A
	papp     uint16 // 'PAPP': '01289',          # puissance Apparente, en VA
	motdetat string // 'MOTDETAT': '000000',     # Mot d'état du compteur
	hhphc    string // 'HHPHC': 'A',             # Horaire Heures Pleines Heures Creuses
	isousc   uint16 // 'ISOUSC': '45',           # Intensité souscrite en A
	hchp     uint64 // 'HCHP': '035972694',      # index heure pleine en Wh
	ptec     string // 'PTEC': 'HP..'            # Période tarifaire en cours
}

// NewLinkyCollector method to construct LinkyCollector
func NewLinkyCollector(device string, baudRate int) *LinkyCollector {
	return &LinkyCollector{
		device:   device,
		baudRate: baudRate,
		optarif: prometheus.NewDesc("linky_optarif",
			"Option tarifaire",
			[]string{"contrat"}, nil,
		),
		imax: prometheus.NewDesc("linky_imax",
			"Intensité max",
			nil, nil,
		),
		hchc: prometheus.NewDesc("linky_hchc",
			"Index heure creuse en Wh",
			nil, nil,
		),
		iinst: prometheus.NewDesc("linky_iinst",
			"Intensité instantanée en A",
			nil, nil,
		),
		papp: prometheus.NewDesc("linky_papp",
			"Puissance Apparente, en VA",
			nil, nil,
		),
		motdetat: prometheus.NewDesc("linky_motdetat",
			"Mot d'état du compteur",
			[]string{"name"}, nil,
		),
		hhphc: prometheus.NewDesc("linky_hhphc",
			"Horaire Heures Pleines Heures Creuses",
			[]string{"name"}, nil,
		),
		isousc: prometheus.NewDesc("linky_isousc",
			"Intensité souscrite en A",
			nil, nil,
		),
		hchp: prometheus.NewDesc("linky_hchp",
			"Index heure pleine en Wh",
			nil, nil,
		),
		ptec: prometheus.NewDesc("linky_ptec",
			"Période tarifaire en cours",
			[]string{"option"}, nil,
		),
	}
}

// Describe implements required describe function for all prometheus collectors
func (collector *LinkyCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.optarif
	ch <- collector.imax
	ch <- collector.hchc
	ch <- collector.iinst
	ch <- collector.papp
	ch <- collector.motdetat
	ch <- collector.hhphc
	ch <- collector.isousc
	ch <- collector.hchp
	ch <- collector.ptec
}

// Collect implements required collect function for all prometheus collectors
func (collector *LinkyCollector) Collect(ch chan<- prometheus.Metric) {
	//for each descriptor or call other functions that do so.
	//Implement logic here to determine proper metric value to return to prometheus
	values := linkyValues{}
	err := collector.readSerial(&values)

	if err == nil {
		//Write latest value for each metric in the prometheus metric channel.
		//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
		ch <- prometheus.MustNewConstMetric(collector.optarif, prometheus.GaugeValue, 1, values.optarif)
		ch <- prometheus.MustNewConstMetric(collector.imax, prometheus.CounterValue, float64(values.imax))
		ch <- prometheus.MustNewConstMetric(collector.hchc, prometheus.CounterValue, float64(values.hchc))
		ch <- prometheus.MustNewConstMetric(collector.iinst, prometheus.CounterValue, float64(values.iinst))
		ch <- prometheus.MustNewConstMetric(collector.papp, prometheus.CounterValue, float64(values.papp))
		ch <- prometheus.MustNewConstMetric(collector.motdetat, prometheus.GaugeValue, 0, values.motdetat)
		ch <- prometheus.MustNewConstMetric(collector.hhphc, prometheus.GaugeValue, 1, values.hhphc)
		ch <- prometheus.MustNewConstMetric(collector.isousc, prometheus.CounterValue, float64(values.isousc))
		ch <- prometheus.MustNewConstMetric(collector.hchp, prometheus.CounterValue, float64(values.hchp))
		switch strings.ToLower(values.ptec) {
		case "hc..":
			ch <- prometheus.MustNewConstMetric(collector.ptec, prometheus.GaugeValue, 1, "hc")
			ch <- prometheus.MustNewConstMetric(collector.ptec, prometheus.GaugeValue, 0, "hp")
		case "hp..":
			ch <- prometheus.MustNewConstMetric(collector.ptec, prometheus.GaugeValue, 0, "hc")
			ch <- prometheus.MustNewConstMetric(collector.ptec, prometheus.GaugeValue, 1, "hp")
		default:
			ch <- prometheus.MustNewConstMetric(collector.ptec, prometheus.GaugeValue, 1, strings.ToLower(values.ptec))
		}
	} else {
		log.Errorf("Unable to read telemetry information : %s", err)
	}
}

// Read information from serial port
func (collector *LinkyCollector) readSerial(linkyValues *linkyValues) error {
	c := &serial.Config{Name: collector.device, Baud: collector.baudRate, Size: 7, Parity: serial.ParityNone, StopBits: serial.Stop1}
	stream, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(stream)
	started := false
	for {
		bytes, _, err := reader.ReadLine()
		if err != nil {
			return err
		}

		line := string(bytes)

		// End loop when block ended
		if started && strings.Contains(line, string(03)) {
			break
		}

		// Start reading data when block started
		if strings.Contains(line, string(02)) {
			started = true
		}

		// Collect data
		if started {
			collector.proceedLine(linkyValues, line)
		}
	}
	return nil
}

// Proceed line by line information
func (collector *LinkyCollector) proceedLine(linkyValues *linkyValues, line string) {
	data := strings.Split(line, " ")

	if len(data) >= 2 {
		name := data[0]
		value := data[1]

		switch strings.ToLower(name) {
		case "optarif":
			linkyValues.optarif = string(value)
		case "imax":
			val, _ := strconv.ParseUint(value, 10, 16)
			linkyValues.imax = uint16(val)
		case "hchc":
			val, _ := strconv.ParseUint(value, 10, 64)
			linkyValues.hchc = val
		case "iinst":
			val, _ := strconv.ParseInt(value, 10, 16)
			linkyValues.iinst = int16(val)
		case "papp":
			val, _ := strconv.ParseUint(value, 10, 16)
			linkyValues.papp = uint16(val)
		case "motdetat":
			linkyValues.motdetat = string(value)
		case "hhphc":
			linkyValues.hhphc = string(value)
		case "isousc":
			val, _ := strconv.ParseUint(value, 10, 16)
			linkyValues.isousc = uint16(val)
		case "hchp":
			val, _ := strconv.ParseUint(value, 10, 64)
			linkyValues.hchp = val
		case "ptec":
			linkyValues.ptec = string(value)
		}
	}
}
