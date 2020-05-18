package collectors

import (
	"fmt"
	"github.com/huin/goserial"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type LinkyCollector struct {
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

type linkyValues struct {
	optarif  string
	imax     uint16
	hchc     uint64
	iinst    int16
	papp     uint16
	motdetat string
	hhphc    string
	isousc   uint16
	hchp     uint64
	ptec     string
}

/*
 'OPTARIF': 'HC..',        # option tarifaire
 'IMAX': '007',            # intensité max
 'HCHC': '040177099',      # index heure creuse en Wh
 'IINST': '005',           # Intensité instantanée en A
 'PAPP': '01289',          # puissance Apparente, en VA
 'MOTDETAT': '000000',     # Mot d'état du compteur
 'HHPHC': 'A',             # Horaire Heures Pleines Heures Creuses
 'ISOUSC': '45',           # Intensité souscrite en A
 'ADCO': '000000000000',   # Adresse du compteur
 'HCHP': '035972694',      # index heure pleine en Wh
 'PTEC': 'HP..'            # Période tarifaire en cours
*/
func NewLinkyCollector() *LinkyCollector {
	return &LinkyCollector{
		optarif: prometheus.NewDesc("linky_optarif",
			"Option tarifaire",
			nil, nil,
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
			nil, nil,
		),
		hhphc: prometheus.NewDesc("linky_hhphc",
			"Horaire Heures Pleines Heures Creuses",
			nil, nil,
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
			nil, nil,
		),
	}
}

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

//Collect implements required collect function for all promehteus collectors
func (collector *LinkyCollector) Collect(ch chan<- prometheus.Metric) {
	//for each descriptor or call other functions that do so.
	//Implement logic here to determine proper metric value to return to prometheus
	values := linkyValues{}
	collector.readSerial(&values)

	//Write latest value for each metric in the prometheus metric channel.
	//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
	ch <- prometheus.MustNewConstMetric(collector.optarif, prometheus.GaugeValue, 1)
	//ch <- prometheus.MustNewConstMetric(collector.imax, prometheus.CounterValue, float64(values.imax))
	//ch <- prometheus.MustNewConstMetric(collector.hchc, prometheus.CounterValue, float64(values.hchc))
	//ch <- prometheus.MustNewConstMetric(collector.iinst, prometheus.CounterValue, float64(values.iinst))
	//ch <- prometheus.MustNewConstMetric(collector.papp, prometheus.CounterValue, float64(values.papp))
	////ch <- prometheus.MustNewConstMetric(collector.motdetat, prometheus.UntypedValue, float64(values.motdetat))
	////ch <- prometheus.MustNewConstMetric(collector.hhphc, prometheus.UntypedValue, float64(values.hhphc))
	//ch <- prometheus.MustNewConstMetric(collector.isousc, prometheus.CounterValue, float64(values.isousc))
	//ch <- prometheus.MustNewConstMetric(collector.hchp, prometheus.CounterValue, float64(values.hchp))
	//ch <- prometheus.MustNewConstMetric(collector.ptec, prometheus.UntypedValue, float64(values.ptec))
}

func (collector *LinkyCollector) readSerial(linkyValues *linkyValues) {
	c := &goserial.Config{Name: "/dev/ttyS0", Baud: 1200, Size: goserial.Byte7, Parity: goserial.ParityNone, StopBits: goserial.StopBits1}
	stream, err := goserial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	buf, err := ioutil.ReadAll(stream)

	if err != nil {
		log.Errorf("Unable to read telemetry information : %s", err)
	}

	fmt.Printf("%s", buf)
}
