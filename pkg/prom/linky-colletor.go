package prom

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/linky-exporter/pkg/core"
)

const USED = "used"
const PRODUCED = "produced"

// LinkyCollector object to describe and collect metrics
type LinkyCollector struct {
	connector              core.LinkyConnector
	linkyDate              *prometheus.Desc
	energyTotal            *prometheus.Desc
	energy                 *prometheus.Desc
	reactiveEnergyTotal    *prometheus.Desc
	intensity              *prometheus.Desc
	voltage                *prometheus.Desc
	power                  *prometheus.Desc
	powerLastYear          *prometheus.Desc
	powerMax               *prometheus.Desc
	powerReference         *prometheus.Desc
	loadCurvePoint         *prometheus.Desc
	loadCurvePointLastYear *prometheus.Desc
	averageVoltage         *prometheus.Desc
	status                 *prometheus.Desc
	movablePeak            *prometheus.Desc
	relay                  *prometheus.Desc
	providerDayInfo        *prometheus.Desc
}

// NewLinkyCollector method to construct LinkyCollector
func NewLinkyCollector(connector core.LinkyConnector) *LinkyCollector {
	return &LinkyCollector{
		connector: connector,
		linkyDate: prometheus.NewDesc("linky_timestamp",
			"Timestamp en seconde",
			[]string{"linky_id", "version", "contract", "pricing"}, nil,
		),
		energyTotal: prometheus.NewDesc("linky_energy_total",
			"Total Energie en Wh",
			[]string{"linky_id", "mode"}, nil,
		),
		energy: prometheus.NewDesc("linky_energy",
			"Energie en Wh",
			[]string{"linky_id", "mode", "index"}, nil,
		),
		reactiveEnergyTotal: prometheus.NewDesc("linky_reactive_energy_total",
			"Total Energie réactive en Wh",
			[]string{"linky_id", "index"}, nil,
		),
		intensity: prometheus.NewDesc("linky_intensity",
			"Courant efficace en A",
			[]string{"linky_id", "phase"}, nil,
		),
		voltage: prometheus.NewDesc("linky_voltage",
			"Tension efficace en V",
			[]string{"linky_id", "phase"}, nil,
		),
		power: prometheus.NewDesc("linky_power",
			"Puissance apparente en VA",
			[]string{"linky_id", "mode", "phase"}, nil,
		),
		powerLastYear: prometheus.NewDesc("linky_power_last_year",
			"Puissance apparente n-1 en VA",
			[]string{"linky_id", "mode", "phase"}, nil,
		),
		powerMax: prometheus.NewDesc("linky_power_max",
			"Puissance apparente en VA",
			[]string{"linky_id", "mode", "phase"}, nil,
		),
		powerReference: prometheus.NewDesc("linky_power_reference",
			"Puissance apparente de référence en kVA",
			[]string{"linky_id", "type"}, nil,
		),
		loadCurvePoint: prometheus.NewDesc("linky_load_curve_point",
			"Point de courbe de charge en W",
			[]string{"linky_id", "mode"}, nil,
		),
		loadCurvePointLastYear: prometheus.NewDesc("linky_load_curve_point_last_year",
			"Point de courbe de charge n-1 en W",
			[]string{"linky_id", "mode"}, nil,
		),
		averageVoltage: prometheus.NewDesc("linky_voltage_average",
			"Tension moyenne en V",
			[]string{"linky_id", "phase"}, nil,
		),
		status: prometheus.NewDesc("linky_status",
			"Statuts issus du registre",
			[]string{"linky_id", "name"}, nil,
		),
		movablePeak: prometheus.NewDesc("linky_movable_peak",
			"Pointe mobile",
			[]string{"linky_id", "type", "phase"}, nil,
		),
		relay: prometheus.NewDesc("linky_relay",
			"Etat du relai",
			[]string{"linky_id", "id"}, nil,
		),
		providerDayInfo: prometheus.NewDesc("linky_provider_day_info",
			"Numéro du jour en cours, du prochain jour et de son profil",
			[]string{"linky_id", "prm", "current_day", "next_day", "next_day_profile"}, nil,
		),
	}
}

// Describe implements required describe function for all prometheus collectors
func (collector *LinkyCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.linkyDate
	ch <- collector.energyTotal
	ch <- collector.energy
	ch <- collector.reactiveEnergyTotal
	ch <- collector.intensity
	ch <- collector.voltage
	ch <- collector.power
	ch <- collector.powerLastYear
	ch <- collector.powerMax
	ch <- collector.powerReference
	ch <- collector.loadCurvePoint
	ch <- collector.loadCurvePointLastYear
	ch <- collector.averageVoltage
	ch <- collector.status
	ch <- collector.movablePeak
	ch <- collector.relay
	ch <- collector.providerDayInfo
}

// Collect implements required collect function for all prometheus collectors
func (collector *LinkyCollector) Collect(ch chan<- prometheus.Metric) {
	var timeSerie LinkyTimeSerie
	var err error

	switch collector.connector.Mode {
	case core.Standard:
		var ticValues *core.StandardTicValue
		ticValues, err = collector.connector.GetLastStandardTicValue()
		if err == nil {
			timeSerie = *ConvertStandardTicValueToTimeSerie(*ticValues)
		}
	case core.Historical:
		var ticValues *core.HistoricalTicValue
		ticValues, err = collector.connector.GetLastHistoricalTicValue()
		if err == nil {
			timeSerie = *ConvertHistoricalTicValueToTimeSerie(*ticValues)
		}
	default:
		log.Error(fmt.Errorf("Not supported mode !"))
	}

	if err == nil {
		// Date
		collector.fillLinkyDateMetric(ch, timeSerie)
		// Energy Total
		collector.fillEnergyTotalMetric(ch, timeSerie)
		// Energy
		collector.fillEnergyMetric(ch, timeSerie)
		// Reactive energy
		collector.fillReactiveEnergyTotalMetric(ch, timeSerie)
		// Intensity
		collector.fillIntensityMetric(ch, timeSerie)
		// Power
		collector.fillPowerMetric(ch, timeSerie)
		// Power Last Year
		collector.fillPowerLastYearMetric(ch, timeSerie)
		// Power Max
		collector.fillPowerMaxMetric(ch, timeSerie)
		// Power Reference
		collector.fillPoweReferenceMetric(ch, timeSerie)
		// Load Curve Point
		collector.fillLoadCurvePointMetric(ch, timeSerie)
		// Load Curve Point Last Year
		collector.fillLoadCurvePointLastYearMetric(ch, timeSerie)
		// Average Voltage
		collector.fillAverageVoltageMetric(ch, timeSerie)

		// Only Standard
		if collector.connector.Mode == core.Standard {
			// Voltage
			collector.fillVoltageMetric(ch, timeSerie)
			// Status
			collector.fillStatusMetric(ch, timeSerie)
			// Relay
			collector.fillRelayMetric(ch, timeSerie)
			// Movable Peak
			if timeSerie.MovingPeakStart1 != 0 {
				collector.fillMovablePeakMetric(ch, timeSerie)
			}
			// Provider Day Info
			// Not enabled now, it's possible to overload metrics cardinality
			// collector.fillProviderDayInfoMetric(ch, timeSerie)
		}
	} else {
		log.Errorf("Unable to read telemetry information : %s", err)
	}
}

// Send to channel linky_date metric
func (collector *LinkyCollector) fillLinkyDateMetric(ch chan<- prometheus.Metric, timeSerie LinkyTimeSerie) {
	ch <- prometheus.MustNewConstMetric(collector.linkyDate, prometheus.CounterValue, timeSerie.LinkyDate, timeSerie.LinkyId, timeSerie.Version, timeSerie.ContractTypeName, timeSerie.PriceLabel)
}

// Send to channel linky_energy_total metric
func (collector *LinkyCollector) fillEnergyTotalMetric(ch chan<- prometheus.Metric, timeSerie LinkyTimeSerie) {
	ch <- prometheus.MustNewConstMetric(collector.energyTotal, prometheus.CounterValue, timeSerie.TotalEnergyUsed, timeSerie.LinkyId, USED)
	if timeSerie.TotalEnergyProduced != 0 {
		ch <- prometheus.MustNewConstMetric(collector.energyTotal, prometheus.CounterValue, timeSerie.TotalEnergyProduced, timeSerie.LinkyId, PRODUCED)
	}
}

// Send to channel linky_energy metric
func (collector *LinkyCollector) fillEnergyMetric(ch chan<- prometheus.Metric, timeSerie LinkyTimeSerie) {
	if timeSerie.EnergyUsedIndex1 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.energy, prometheus.CounterValue, timeSerie.EnergyUsedIndex1, timeSerie.LinkyId, USED, "F1")
	}
	if timeSerie.EnergyUsedIndex2 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.energy, prometheus.CounterValue, timeSerie.EnergyUsedIndex2, timeSerie.LinkyId, USED, "F2")
	}
	if timeSerie.EnergyUsedIndex3 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.energy, prometheus.CounterValue, timeSerie.EnergyUsedIndex3, timeSerie.LinkyId, USED, "F3")
	}
	if timeSerie.EnergyUsedIndex4 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.energy, prometheus.CounterValue, timeSerie.EnergyUsedIndex4, timeSerie.LinkyId, USED, "F4")
	}
	if timeSerie.EnergyUsedIndex5 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.energy, prometheus.CounterValue, timeSerie.EnergyUsedIndex5, timeSerie.LinkyId, USED, "F5")
	}
	if timeSerie.EnergyUsedIndex6 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.energy, prometheus.CounterValue, timeSerie.EnergyUsedIndex6, timeSerie.LinkyId, USED, "F6")
	}
	if timeSerie.EnergyUsedIndex7 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.energy, prometheus.CounterValue, timeSerie.EnergyUsedIndex7, timeSerie.LinkyId, USED, "F7")
	}
	if timeSerie.EnergyUsedIndex8 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.energy, prometheus.CounterValue, timeSerie.EnergyUsedIndex8, timeSerie.LinkyId, USED, "F8")
	}
	if timeSerie.EnergyUsedIndex9 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.energy, prometheus.CounterValue, timeSerie.EnergyUsedIndex9, timeSerie.LinkyId, USED, "F9")
	}
	if timeSerie.EnergyUsedIndex10 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.energy, prometheus.CounterValue, timeSerie.EnergyUsedIndex10, timeSerie.LinkyId, USED, "F10")
	}
	if timeSerie.EnergyUsedDistributorIndex1 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.energy, prometheus.CounterValue, timeSerie.EnergyUsedDistributorIndex1, timeSerie.LinkyId, USED, "D1")
	}
	if timeSerie.EnergyUsedDistributorIndex2 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.energy, prometheus.CounterValue, timeSerie.EnergyUsedDistributorIndex2, timeSerie.LinkyId, USED, "D2")
	}
	if timeSerie.EnergyUsedDistributorIndex3 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.energy, prometheus.CounterValue, timeSerie.EnergyUsedDistributorIndex3, timeSerie.LinkyId, USED, "D3")
	}
	if timeSerie.EnergyUsedDistributorIndex4 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.energy, prometheus.CounterValue, timeSerie.EnergyUsedDistributorIndex4, timeSerie.LinkyId, USED, "D4")
	}
}

// Send to channel linky_reactive_energy_total metric
func (collector *LinkyCollector) fillReactiveEnergyTotalMetric(ch chan<- prometheus.Metric, timeSerie LinkyTimeSerie) {
	if timeSerie.TotalReactiveEnergyQ1 != 0 || timeSerie.TotalReactiveEnergyQ2 != 0 || timeSerie.TotalReactiveEnergyQ3 != 0 || timeSerie.TotalReactiveEnergyQ4 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.reactiveEnergyTotal, prometheus.CounterValue, timeSerie.TotalReactiveEnergyQ1, timeSerie.LinkyId, "Q1")
		ch <- prometheus.MustNewConstMetric(collector.reactiveEnergyTotal, prometheus.CounterValue, timeSerie.TotalReactiveEnergyQ2, timeSerie.LinkyId, "Q2")
		ch <- prometheus.MustNewConstMetric(collector.reactiveEnergyTotal, prometheus.CounterValue, timeSerie.TotalReactiveEnergyQ3, timeSerie.LinkyId, "Q3")
		ch <- prometheus.MustNewConstMetric(collector.reactiveEnergyTotal, prometheus.CounterValue, timeSerie.TotalReactiveEnergyQ4, timeSerie.LinkyId, "Q4")
	}
}

// Send to channel linky_intensity metric
func (collector *LinkyCollector) fillIntensityMetric(ch chan<- prometheus.Metric, timeSerie LinkyTimeSerie) {
	ch <- prometheus.MustNewConstMetric(collector.intensity, prometheus.GaugeValue, timeSerie.IntensityP1, timeSerie.LinkyId, "1")
	if timeSerie.IntensityP2 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.intensity, prometheus.GaugeValue, timeSerie.IntensityP2, timeSerie.LinkyId, "2")
	}
	if timeSerie.IntensityP3 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.intensity, prometheus.GaugeValue, timeSerie.IntensityP3, timeSerie.LinkyId, "3")
	}
}

// Send to channel linky_voltage metric
func (collector *LinkyCollector) fillVoltageMetric(ch chan<- prometheus.Metric, timeSerie LinkyTimeSerie) {
	ch <- prometheus.MustNewConstMetric(collector.voltage, prometheus.GaugeValue, timeSerie.VoltageP1, timeSerie.LinkyId, "1")
	if timeSerie.VoltageP2 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.voltage, prometheus.GaugeValue, timeSerie.VoltageP1, timeSerie.LinkyId, "2")
	}
	if timeSerie.VoltageP3 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.voltage, prometheus.GaugeValue, timeSerie.VoltageP1, timeSerie.LinkyId, "3")
	}
}

// Send to channel linky_power metric
func (collector *LinkyCollector) fillPowerMetric(ch chan<- prometheus.Metric, timeSerie LinkyTimeSerie) {
	if timeSerie.PowerUsed != 0 {
		ch <- prometheus.MustNewConstMetric(collector.power, prometheus.GaugeValue, timeSerie.PowerUsed, timeSerie.LinkyId, USED, "1")
	}
	if timeSerie.PowerUsedP1 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.power, prometheus.GaugeValue, timeSerie.PowerUsedP1, timeSerie.LinkyId, USED, "1")
	}
	if timeSerie.PowerUsedP2 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.power, prometheus.GaugeValue, timeSerie.PowerUsedP2, timeSerie.LinkyId, USED, "2")
	}
	if timeSerie.PowerUsedP3 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.power, prometheus.GaugeValue, timeSerie.PowerUsedP3, timeSerie.LinkyId, USED, "3")
	}
	if timeSerie.PowerProduced != 0 {
		ch <- prometheus.MustNewConstMetric(collector.power, prometheus.GaugeValue, timeSerie.PowerProduced, timeSerie.LinkyId, PRODUCED, "0")
	}
}

// Send to channel linky_power_last_year metric
func (collector *LinkyCollector) fillPowerLastYearMetric(ch chan<- prometheus.Metric, timeSerie LinkyTimeSerie) {
	if timeSerie.PowerUsedMaxLastYear != 0 {
		ch <- prometheus.MustNewConstMetric(collector.powerLastYear, prometheus.GaugeValue, timeSerie.PowerUsedMaxLastYear, timeSerie.LinkyId, USED, "1")
	}
	if timeSerie.PowerUsedMaxLastYearP1 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.powerLastYear, prometheus.GaugeValue, timeSerie.PowerUsedMaxLastYearP1, timeSerie.LinkyId, USED, "1")
	}
	if timeSerie.PowerUsedMaxLastYearP2 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.powerLastYear, prometheus.GaugeValue, timeSerie.PowerUsedMaxLastYearP2, timeSerie.LinkyId, USED, "2")
	}
	if timeSerie.PowerUsedMaxLastYearP3 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.powerLastYear, prometheus.GaugeValue, timeSerie.PowerUsedMaxLastYearP3, timeSerie.LinkyId, USED, "3")
	}
	if timeSerie.PowerProducedLastYear != 0 {
		ch <- prometheus.MustNewConstMetric(collector.powerLastYear, prometheus.GaugeValue, timeSerie.PowerProducedLastYear, timeSerie.LinkyId, PRODUCED, "0")
	}
}

// Send to channel linky_power_max metric
func (collector *LinkyCollector) fillPowerMaxMetric(ch chan<- prometheus.Metric, timeSerie LinkyTimeSerie) {
	if timeSerie.PowerUsedMax != 0 {
		ch <- prometheus.MustNewConstMetric(collector.powerMax, prometheus.GaugeValue, timeSerie.PowerUsedMax, timeSerie.LinkyId, USED, "1")
	}
	if timeSerie.PowerUsedMaxP1 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.powerMax, prometheus.GaugeValue, timeSerie.PowerUsedMaxP1, timeSerie.LinkyId, USED, "1")
	}
	if timeSerie.PowerUsedMaxP2 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.powerMax, prometheus.GaugeValue, timeSerie.PowerUsedMaxP2, timeSerie.LinkyId, USED, "2")
	}
	if timeSerie.PowerUsedMaxP3 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.powerMax, prometheus.GaugeValue, timeSerie.PowerUsedMaxP3, timeSerie.LinkyId, USED, "3")
	}
	if timeSerie.PowerProducedMax != 0 {
		ch <- prometheus.MustNewConstMetric(collector.powerMax, prometheus.GaugeValue, timeSerie.PowerProducedMax, timeSerie.LinkyId, PRODUCED, "0")
	}
}

// Send to channel linky_power_reference metric
func (collector *LinkyCollector) fillPoweReferenceMetric(ch chan<- prometheus.Metric, timeSerie LinkyTimeSerie) {
	if timeSerie.ReferencePower != 0 {
		ch <- prometheus.MustNewConstMetric(collector.powerReference, prometheus.GaugeValue, timeSerie.ReferencePower, timeSerie.LinkyId, "subscribed")
	}
	if timeSerie.BreakingPower != 0 {
		ch <- prometheus.MustNewConstMetric(collector.powerReference, prometheus.GaugeValue, timeSerie.BreakingPower, timeSerie.LinkyId, "breaking")
	}
}

// Send to channel linky_load_curve_point metric
func (collector *LinkyCollector) fillLoadCurvePointMetric(ch chan<- prometheus.Metric, timeSerie LinkyTimeSerie) {
	if timeSerie.UsedLoadCurvePoint != 0 {
		ch <- prometheus.MustNewConstMetric(collector.loadCurvePoint, prometheus.GaugeValue, timeSerie.UsedLoadCurvePoint, timeSerie.LinkyId, USED)
	}
	if timeSerie.ProducedLoadCurvePoint != 0 {
		ch <- prometheus.MustNewConstMetric(collector.loadCurvePoint, prometheus.GaugeValue, timeSerie.ProducedLoadCurvePoint, timeSerie.LinkyId, PRODUCED)
	}
}

// Send to channel linky_load_curve_point_last_year metric
func (collector *LinkyCollector) fillLoadCurvePointLastYearMetric(ch chan<- prometheus.Metric, timeSerie LinkyTimeSerie) {
	if timeSerie.UsedLoadCurvePoint != 0 {
		ch <- prometheus.MustNewConstMetric(collector.loadCurvePointLastYear, prometheus.GaugeValue, timeSerie.UsedLoadCurvePointLastYear, timeSerie.LinkyId, USED)
	}
	if timeSerie.ProducedLoadCurvePoint != 0 {
		ch <- prometheus.MustNewConstMetric(collector.loadCurvePointLastYear, prometheus.GaugeValue, timeSerie.ProducedLoadCurvePointLastYear, timeSerie.LinkyId, PRODUCED)
	}
}

// Send to channel linky_average_voltage metric
func (collector *LinkyCollector) fillAverageVoltageMetric(ch chan<- prometheus.Metric, timeSerie LinkyTimeSerie) {
	if timeSerie.AverageVoltageP1 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.averageVoltage, prometheus.GaugeValue, timeSerie.AverageVoltageP1, timeSerie.LinkyId, "1")
	}
	if timeSerie.AverageVoltageP2 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.averageVoltage, prometheus.GaugeValue, timeSerie.AverageVoltageP2, timeSerie.LinkyId, "2")
	}
	if timeSerie.AverageVoltageP3 != 0 {
		ch <- prometheus.MustNewConstMetric(collector.averageVoltage, prometheus.GaugeValue, timeSerie.AverageVoltageP3, timeSerie.LinkyId, "3")
	}
}

// Send to channel linky_status metric
func (collector *LinkyCollector) fillStatusMetric(ch chan<- prometheus.Metric, timeSerie LinkyTimeSerie) {
	ch <- prometheus.MustNewConstMetric(collector.status, prometheus.GaugeValue, timeSerie.DryContactStatus, timeSerie.LinkyId, "Contact sec")
	ch <- prometheus.MustNewConstMetric(collector.status, prometheus.GaugeValue, timeSerie.CutOffDeviceStatus, timeSerie.LinkyId, "Organe de coupure")
	ch <- prometheus.MustNewConstMetric(collector.status, prometheus.GaugeValue, timeSerie.LinkyTerminalShieldStatus, timeSerie.LinkyId, "État du cache-bornes distributeur")
	ch <- prometheus.MustNewConstMetric(collector.status, prometheus.GaugeValue, timeSerie.SurgeStatus, timeSerie.LinkyId, "Surtension sur une des phases")
	ch <- prometheus.MustNewConstMetric(collector.status, prometheus.GaugeValue, timeSerie.ReferencePowerExceededStatus, timeSerie.LinkyId, "Dépassement de la puissance de référence")
	ch <- prometheus.MustNewConstMetric(collector.status, prometheus.GaugeValue, timeSerie.ConsumptionStatus, timeSerie.LinkyId, "Fonctionnement producteur/consommateur")
	ch <- prometheus.MustNewConstMetric(collector.status, prometheus.GaugeValue, timeSerie.EnergyDirectionStatus, timeSerie.LinkyId, "Sens de l énergie active")
	ch <- prometheus.MustNewConstMetric(collector.status, prometheus.GaugeValue, timeSerie.ContractTypePriceStatus, timeSerie.LinkyId, "Tarif en cours sur le contrat fourniture")
	ch <- prometheus.MustNewConstMetric(collector.status, prometheus.GaugeValue, timeSerie.ContractTypePriceDistributorStatus, timeSerie.LinkyId, "Tarif en cours sur le contrat distributeur")
	ch <- prometheus.MustNewConstMetric(collector.status, prometheus.GaugeValue, timeSerie.ClockStatus, timeSerie.LinkyId, "Mode dégradée de l horloge")
	ch <- prometheus.MustNewConstMetric(collector.status, prometheus.GaugeValue, timeSerie.TicStatus, timeSerie.LinkyId, "État de la sortie télé-information")
	ch <- prometheus.MustNewConstMetric(collector.status, prometheus.GaugeValue, timeSerie.EuridisLinkStatus, timeSerie.LinkyId, "État de la sortie communication Euridis")
	ch <- prometheus.MustNewConstMetric(collector.status, prometheus.GaugeValue, timeSerie.CPLStatus, timeSerie.LinkyId, "Statut du CPL")
	ch <- prometheus.MustNewConstMetric(collector.status, prometheus.GaugeValue, timeSerie.CPLSyncStatus, timeSerie.LinkyId, "Synchronisation CPL")
	ch <- prometheus.MustNewConstMetric(collector.status, prometheus.GaugeValue, timeSerie.TempoContractColorStatus, timeSerie.LinkyId, "Couleur du jour pour le contrat historique tempo")
	ch <- prometheus.MustNewConstMetric(collector.status, prometheus.GaugeValue, timeSerie.TempoContractNextDayColorStatus, timeSerie.LinkyId, "Couleur du lendemain pour le contrat historique tempo")
	ch <- prometheus.MustNewConstMetric(collector.status, prometheus.GaugeValue, timeSerie.MovingPeakNoticeStatus, timeSerie.LinkyId, "Préavis pointes mobiles")
	ch <- prometheus.MustNewConstMetric(collector.status, prometheus.GaugeValue, timeSerie.MovingPeakStatus, timeSerie.LinkyId, "Pointe mobile (PM)")
}

// Send to channel linky_movable_peak metric
func (collector *LinkyCollector) fillMovablePeakMetric(ch chan<- prometheus.Metric, timeSerie LinkyTimeSerie) {
	ch <- prometheus.MustNewConstMetric(collector.movablePeak, prometheus.GaugeValue, timeSerie.MovingPeakStart1, timeSerie.LinkyId, "start", "1")
	ch <- prometheus.MustNewConstMetric(collector.movablePeak, prometheus.GaugeValue, timeSerie.MovingPeakEnd1, timeSerie.LinkyId, "end", "1")
	ch <- prometheus.MustNewConstMetric(collector.movablePeak, prometheus.GaugeValue, timeSerie.MovingPeakStart2, timeSerie.LinkyId, "start", "2")
	ch <- prometheus.MustNewConstMetric(collector.movablePeak, prometheus.GaugeValue, timeSerie.MovingPeakEnd2, timeSerie.LinkyId, "end", "2")
	ch <- prometheus.MustNewConstMetric(collector.movablePeak, prometheus.GaugeValue, timeSerie.MovingPeakStart3, timeSerie.LinkyId, "start", "3")
	ch <- prometheus.MustNewConstMetric(collector.movablePeak, prometheus.GaugeValue, timeSerie.MovingPeakEnd3, timeSerie.LinkyId, "end", "3")
}

// Send to channel linky_relay metric
func (collector *LinkyCollector) fillRelayMetric(ch chan<- prometheus.Metric, timeSerie LinkyTimeSerie) {
	ch <- prometheus.MustNewConstMetric(collector.relay, prometheus.GaugeValue, timeSerie.Relay1, timeSerie.LinkyId, "1")
	ch <- prometheus.MustNewConstMetric(collector.relay, prometheus.GaugeValue, timeSerie.Relay2, timeSerie.LinkyId, "2")
	ch <- prometheus.MustNewConstMetric(collector.relay, prometheus.GaugeValue, timeSerie.Relay3, timeSerie.LinkyId, "3")
	ch <- prometheus.MustNewConstMetric(collector.relay, prometheus.GaugeValue, timeSerie.Relay4, timeSerie.LinkyId, "4")
	ch <- prometheus.MustNewConstMetric(collector.relay, prometheus.GaugeValue, timeSerie.Relay5, timeSerie.LinkyId, "5")
	ch <- prometheus.MustNewConstMetric(collector.relay, prometheus.GaugeValue, timeSerie.Relay6, timeSerie.LinkyId, "6")
	ch <- prometheus.MustNewConstMetric(collector.relay, prometheus.GaugeValue, timeSerie.Relay7, timeSerie.LinkyId, "7")
	ch <- prometheus.MustNewConstMetric(collector.relay, prometheus.GaugeValue, timeSerie.Relay8, timeSerie.LinkyId, "8")
}

// Send to channel linky_provider_day_info metric
func (collector *LinkyCollector) fillProviderDayInfoMetric(ch chan<- prometheus.Metric, timeSerie LinkyTimeSerie) {
	ch <- prometheus.MustNewConstMetric(collector.providerDayInfo, prometheus.GaugeValue, 1, timeSerie.LinkyId, timeSerie.Prm, timeSerie.ContractTypeDayNumber, timeSerie.ContractTypeNextDayNumber, timeSerie.ContractTypeNextDayProfile)
}
