package collectors

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

// LinkyCollector object to describe and collect metrics
type LinkyCollector struct {
	device              string
	baudRate            int
	frameSize           byte
	parity              serial.Parity
	stopBits            serial.StopBits
	index               *prometheus.Desc
	power               *prometheus.Desc
	intensity           *prometheus.Desc
	intensitySubscribed *prometheus.Desc
	intensityMax        *prometheus.Desc
	tomorrowBlue        *prometheus.Desc
	tomorrowWhite       *prometheus.Desc
	tomorrowRed         *prometheus.Desc
	hoursGroup          *prometheus.Desc
}

// Internal linky values object to each metrics
type linkyValues struct {
	adco    string // 'ADCO': '2000..',         # Identification de compteur
	optarif string // 'OPTARIF': 'HC..',        # option tarifaire
	base    uint64 // 'BASE': '040177099',      # index tarif de base
	imax    uint16 // 'IMAX': '007',            # intensité max
	hchc    uint64 // 'HCHC': '040177099',      # index heure creuse en Wh
	iinst   int16  // 'IINST': '005',           # Intensité instantanée en A
	papp    uint16 // 'PAPP': '01289',          # puissance Apparente, en VA
	hhphc   string // 'HHPHC': 'A',             # Horaire Heures Pleines Heures Creuses
	isousc  uint16 // 'ISOUSC': '45',           # Intensité souscrite en A
	hchp    uint64 // 'HCHP': '040177099',      # index heure pleine en Wh
	ptec    string // 'PTEC': 'HP..'            # Période tarifaire en cours
	hpjb    uint64 // 'HPJB': '040177099',      # index heures creuses jours bleus en wh
	hcjb    uint64 // 'HCJB': '040177099',      # index heures pleines jours bleus en wh
	hpjw    uint64 // 'HPJW': '040177099',      # index heures creuses jours blancs en wh
	hcjw    uint64 // 'HCJW': '040177099',      # index heures pleines jours blancs en wh
	hpjr    uint64 // 'HPJR': '040177099',      # index heures creuses jours rouges en wh
	hcjr    uint64 // 'HCJR': '040177099',      # index heures pleines jours rouges en wh
	demain  string // 'DEMAIN': 'BLAN'          # Couleur du lendemain
}

// NewLinkyCollector method to construct LinkyCollector
func NewLinkyCollector(device string, baudRate int, frameSize byte, parity serial.Parity, stopBits serial.StopBits) *LinkyCollector {
	return &LinkyCollector{
		device:    device,
		baudRate:  baudRate,
		frameSize: frameSize,
		parity:    parity,
		stopBits:  stopBits,
		index: prometheus.NewDesc("linky_index_watthours_total",
			"Index en Wh",
			[]string{"idcompteur", "tarif", "periode"}, nil,
		),
		power: prometheus.NewDesc("linky_power_voltamperes",
			"Puissance apparente en VA",
			[]string{"idcompteur", "tarif"}, nil,
		),
		intensity: prometheus.NewDesc("linky_intensity_amperes",
			"Intensité en A",
			[]string{"idcompteur", "tarif"}, nil,
		),
		intensitySubscribed: prometheus.NewDesc("linky_subscribed_intensity_amperes",
			"Intensité souscrite en A",
			[]string{"idcompteur", "tarif"}, nil,
		),
		intensityMax: prometheus.NewDesc("linky_maximum_intensity_amperes",
			"Intensité maximale en A",
			[]string{"idcompteur", "tarif"}, nil,
		),
		tomorrowBlue: prometheus.NewDesc("linky_tomorrow_blue_info",
			"Lendemain Tempo bleu",
			[]string{"idcompteur", "tarif"}, nil,
		),
		tomorrowWhite: prometheus.NewDesc("linky_tomorrow_white_info",
			"Lendemain Tempo blanc",
			[]string{"idcompteur", "tarif"}, nil,
		),
		tomorrowRed: prometheus.NewDesc("linky_tomorrow_red_info",
			"Lendemain Tempo rouge",
			[]string{"idcompteur", "tarif"}, nil,
		),
		hoursGroup: prometheus.NewDesc("linky_hours_group_info",
			"Groupe horaire (tarif Tempo ou HPHC)",
			[]string{"idcompteur", "tarif", "groupe"}, nil,
		),
	}
}

// Describe implements required describe function for all prometheus collectors
func (collector *LinkyCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.index
	ch <- collector.power
	ch <- collector.intensity
	ch <- collector.intensitySubscribed
	ch <- collector.intensityMax
	ch <- collector.tomorrowBlue
	ch <- collector.tomorrowWhite
	ch <- collector.tomorrowRed
	ch <- collector.hoursGroup
}

// Collect implements required collect function for all prometheus collectors
func (collector *LinkyCollector) Collect(ch chan<- prometheus.Metric) {
	//for each descriptor or call other functions that do so.
	//Implement logic here to determine proper metric value to return to prometheus
	values := linkyValues{}
	err := collector.readSerial(&values)
	var tarif string

	if err == nil {
		switch strings.ToLower(values.optarif) {
		case "base":
			tarif = "base"
		case "hc..":
			tarif = "heures creuses"
		case "bbrx":
			tarif = "tempo"
		default:
			tarif = values.optarif
		}
		//Write latest value for each metric in the prometheus metric channel.
		//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
		ch <- prometheus.MustNewConstMetric(collector.power, prometheus.GaugeValue, float64(values.papp), values.adco, tarif)
		ch <- prometheus.MustNewConstMetric(collector.intensity, prometheus.GaugeValue, float64(values.iinst), values.adco, tarif)
		ch <- prometheus.MustNewConstMetric(collector.intensitySubscribed, prometheus.GaugeValue, float64(values.isousc), values.adco, tarif)
		ch <- prometheus.MustNewConstMetric(collector.intensityMax, prometheus.GaugeValue, float64(values.imax), values.adco, tarif)
		ch <- prometheus.MustNewConstMetric(collector.hoursGroup, prometheus.GaugeValue, 1, values.adco, tarif, values.hhphc)
		switch strings.ToLower(values.ptec) {
		case "th..":
			ch <- prometheus.MustNewConstMetric(collector.index, prometheus.CounterValue, float64(values.base), values.adco, tarif, "-")
		case "hc..":
			ch <- prometheus.MustNewConstMetric(collector.index, prometheus.CounterValue, float64(values.hchc), values.adco, tarif, "HC")
		case "hp..":
			ch <- prometheus.MustNewConstMetric(collector.index, prometheus.CounterValue, float64(values.hchp), values.adco, tarif, "HP")
		case "hcjb":
			ch <- prometheus.MustNewConstMetric(collector.index, prometheus.CounterValue, float64(values.hcjb), values.adco, tarif, "HCJB")
		case "hcjw":
			ch <- prometheus.MustNewConstMetric(collector.index, prometheus.CounterValue, float64(values.hcjw), values.adco, tarif, "HCJW")
		case "hcjr":
			ch <- prometheus.MustNewConstMetric(collector.index, prometheus.CounterValue, float64(values.hcjr), values.adco, tarif, "HCJR")
		case "hpjb":
			ch <- prometheus.MustNewConstMetric(collector.index, prometheus.CounterValue, float64(values.hpjb), values.adco, tarif, "HPJB")
		case "hpjw":
			ch <- prometheus.MustNewConstMetric(collector.index, prometheus.CounterValue, float64(values.hpjw), values.adco, tarif, "HPJW")
		case "hpjr":
			ch <- prometheus.MustNewConstMetric(collector.index, prometheus.CounterValue, float64(values.hpjr), values.adco, tarif, "HPJR")
		default:
		}
		switch strings.ToLower(values.demain) {
		case "bleu":
			ch <- prometheus.MustNewConstMetric(collector.tomorrowBlue, prometheus.GaugeValue, 1, values.adco, tarif)
			ch <- prometheus.MustNewConstMetric(collector.tomorrowWhite, prometheus.GaugeValue, 0, values.adco, tarif)
			ch <- prometheus.MustNewConstMetric(collector.tomorrowRed, prometheus.GaugeValue, 0, values.adco, tarif)
		case "blan":
			ch <- prometheus.MustNewConstMetric(collector.tomorrowBlue, prometheus.GaugeValue, 0, values.adco, tarif)
			ch <- prometheus.MustNewConstMetric(collector.tomorrowWhite, prometheus.GaugeValue, 1, values.adco, tarif)
			ch <- prometheus.MustNewConstMetric(collector.tomorrowRed, prometheus.GaugeValue, 0, values.adco, tarif)
		case "roug":
			ch <- prometheus.MustNewConstMetric(collector.tomorrowBlue, prometheus.GaugeValue, 0, values.adco, tarif)
			ch <- prometheus.MustNewConstMetric(collector.tomorrowWhite, prometheus.GaugeValue, 0, values.adco, tarif)
			ch <- prometheus.MustNewConstMetric(collector.tomorrowRed, prometheus.GaugeValue, 1, values.adco, tarif)
		}
	} else {
		log.Errorf("Unable to read telemetry information : %s", err)
	}
}

// Read information from serial port
func (collector *LinkyCollector) readSerial(linkyValues *linkyValues) error {
	c := &serial.Config{Name: collector.device, Baud: collector.baudRate, Size: collector.frameSize, Parity: collector.parity, StopBits: collector.stopBits}
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
		case "adco":
			linkyValues.adco = string(value)
		case "optarif":
			linkyValues.optarif = string(value)
		case "base":
			val, _ := strconv.ParseUint(value, 10, 64)
			linkyValues.base = val
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
		case "hhphc":
			linkyValues.hhphc = string(value)
		case "isousc":
			val, _ := strconv.ParseUint(value, 10, 16)
			linkyValues.isousc = uint16(val)
		case "hchp":
			val, _ := strconv.ParseUint(value, 10, 64)
			linkyValues.hchp = val
		case "bbrhcjb":
			val, _ := strconv.ParseUint(value, 10, 64)
			linkyValues.hcjb = val
		case "bbrhpjb":
			val, _ := strconv.ParseUint(value, 10, 64)
			linkyValues.hpjb = val
		case "bbrhcjw":
			val, _ := strconv.ParseUint(value, 10, 64)
			linkyValues.hcjw = val
		case "bbrhpjw":
			val, _ := strconv.ParseUint(value, 10, 64)
			linkyValues.hpjw = val
		case "bbrhcjr":
			val, _ := strconv.ParseUint(value, 10, 64)
			linkyValues.hcjr = val
		case "bbrhpjr":
			val, _ := strconv.ParseUint(value, 10, 64)
			linkyValues.hpjr = val
		case "demain":
			linkyValues.demain = string(value)
		case "ptec":
			linkyValues.ptec = string(value)
		}
	}
}
