package prom

import (
	"strconv"

	"github.com/syberalexis/linky-exporter/pkg/core"
)

// Convert (with construction) Historical Tic Value to Time serie value
func ConvertHistoricalTicValueToTimeSerie(historicalValues core.HistoricalTicValue) *LinkyTimeSerie {
	timeSerie := &LinkyTimeSerie{
		LinkyId:          historicalValues.Adco,
		Version:          "1",
		ContractTypeName: historicalValues.Optarif,
		PriceLabel:       historicalValues.Ptec,
		PowerUsed:        float64(historicalValues.Papp),
	}

	isTriplePhase := historicalValues.Iinst2 != 0 || historicalValues.Iinst3 != 0
	isBase := historicalValues.Base != 0
	isHCHP := historicalValues.Hchc != 0 || historicalValues.Hchp != 0
	isEJP := historicalValues.Ejphn != 0 || historicalValues.Ejphpn != 0
	isBBR := historicalValues.Bbrhcjb != 0 || historicalValues.Bbrhpjb != 0 ||
		historicalValues.Bbrhcjw != 0 || historicalValues.Bbrhpjw != 0 ||
		historicalValues.Bbrhcjr != 0 || historicalValues.Bbrhpjr != 0

	if isTriplePhase {
		timeSerie.ReferencePower = float64(historicalValues.Isousc) * 3 * 200 / 1000
		timeSerie.IntensityP1 = float64(historicalValues.Iinst1)
		timeSerie.IntensityP2 = float64(historicalValues.Iinst2)
		timeSerie.IntensityP3 = float64(historicalValues.Iinst3)
	} else {
		timeSerie.ReferencePower = float64(historicalValues.Isousc) * 200 / 1000
		timeSerie.IntensityP1 = float64(historicalValues.Iinst)
		timeSerie.BreakingPower = float64(historicalValues.Adps) * 200 / 1000
	}

	if isBase {
		timeSerie.EnergyUsedIndex1 = float64(historicalValues.Base)
	} else if isHCHP {
		timeSerie.EnergyUsedIndex1 = float64(historicalValues.Hchc)
		timeSerie.EnergyUsedIndex2 = float64(historicalValues.Hchp)
		timeSerie.TotalEnergyUsed = float64(historicalValues.Hchc + historicalValues.Hchp)
		timeSerie.ContractTypeDayNumber = historicalValues.Hhphc
	} else if isEJP {
		timeSerie.EnergyUsedIndex1 = float64(historicalValues.Ejphn)
		timeSerie.EnergyUsedIndex2 = float64(historicalValues.Ejphpn)
		timeSerie.ContractTypeNextDayNumber = strconv.FormatInt(int64(historicalValues.Pejp), 10)
	} else if isBBR {
		timeSerie.EnergyUsedIndex1 = float64(historicalValues.Bbrhcjb)
		timeSerie.EnergyUsedIndex2 = float64(historicalValues.Bbrhpjb)
		timeSerie.EnergyUsedIndex3 = float64(historicalValues.Bbrhcjw)
		timeSerie.EnergyUsedIndex4 = float64(historicalValues.Bbrhpjw)
		timeSerie.EnergyUsedIndex5 = float64(historicalValues.Bbrhcjr)
		timeSerie.EnergyUsedIndex6 = float64(historicalValues.Bbrhpjr)
		timeSerie.ContractTypeNextDayNumber = historicalValues.Demain
	}

	return timeSerie
}

// Convert Standard Tic Value to Time serie value
func ConvertStandardTicValueToTimeSerie(standardValues core.StandardTicValue) *LinkyTimeSerie {
	return &LinkyTimeSerie{
		LinkyId:                            standardValues.Adsc,
		Version:                            standardValues.Vtic,
		LinkyDate:                          float64(standardValues.Date.Unix()),
		ContractTypeName:                   standardValues.Ngtf,
		PriceLabel:                         standardValues.Ltarf,
		TotalEnergyUsed:                    float64(standardValues.East),
		EnergyUsedIndex1:                   float64(standardValues.Easf01),
		EnergyUsedIndex2:                   float64(standardValues.Easf02),
		EnergyUsedIndex3:                   float64(standardValues.Easf03),
		EnergyUsedIndex4:                   float64(standardValues.Easf04),
		EnergyUsedIndex5:                   float64(standardValues.Easf05),
		EnergyUsedIndex6:                   float64(standardValues.Easf06),
		EnergyUsedIndex7:                   float64(standardValues.Easf07),
		EnergyUsedIndex8:                   float64(standardValues.Easf08),
		EnergyUsedIndex9:                   float64(standardValues.Easf09),
		EnergyUsedIndex10:                  float64(standardValues.Easf10),
		EnergyUsedDistributorIndex1:        float64(standardValues.Easd01),
		EnergyUsedDistributorIndex2:        float64(standardValues.Easd02),
		EnergyUsedDistributorIndex3:        float64(standardValues.Easd03),
		EnergyUsedDistributorIndex4:        float64(standardValues.Easd04),
		TotalEnergyProduced:                float64(standardValues.Eait),
		TotalReactiveEnergyQ1:              float64(standardValues.Erq1),
		TotalReactiveEnergyQ2:              float64(standardValues.Erq2),
		TotalReactiveEnergyQ3:              float64(standardValues.Erq3),
		TotalReactiveEnergyQ4:              float64(standardValues.Erq4),
		IntensityP1:                        float64(standardValues.Irms1),
		IntensityP2:                        float64(standardValues.Irms2),
		IntensityP3:                        float64(standardValues.Irms3),
		VoltageP1:                          float64(standardValues.Urms1),
		VoltageP2:                          float64(standardValues.Urms2),
		VoltageP3:                          float64(standardValues.Urms3),
		ReferencePower:                     float64(standardValues.Pref),
		BreakingPower:                      float64(standardValues.Pcoup),
		PowerUsed:                          float64(standardValues.Sinsts),
		PowerUsedP1:                        float64(standardValues.Sinsts1),
		PowerUsedP2:                        float64(standardValues.Sinsts2),
		PowerUsedP3:                        float64(standardValues.Sinsts3),
		PowerUsedMax:                       float64(standardValues.Smaxsn),
		PowerUsedMaxP1:                     float64(standardValues.Smaxsn1),
		PowerUsedMaxP2:                     float64(standardValues.Smaxsn2),
		PowerUsedMaxP3:                     float64(standardValues.Smaxsn3),
		PowerUsedMaxLastYear:               float64(standardValues.Smaxsnly),
		PowerUsedMaxLastYearP1:             float64(standardValues.Smaxsn1ly),
		PowerUsedMaxLastYearP2:             float64(standardValues.Smaxsn2ly),
		PowerUsedMaxLastYearP3:             float64(standardValues.Smaxsn3ly),
		PowerProduced:                      float64(standardValues.Sinsti),
		PowerProducedMax:                   float64(standardValues.Smaxin),
		PowerProducedLastYear:              float64(standardValues.Smaxinly),
		UsedLoadCurvePoint:                 float64(standardValues.Ccasn),
		UsedLoadCurvePointLastYear:         float64(standardValues.Ccasnly),
		ProducedLoadCurvePoint:             float64(standardValues.Ccain),
		ProducedLoadCurvePointLastYear:     float64(standardValues.Ccainly),
		AverageVoltageP1:                   float64(standardValues.Umoy1),
		AverageVoltageP2:                   float64(standardValues.Umoy2),
		AverageVoltageP3:                   float64(standardValues.Umoy3),
		DryContactStatus:                   float64(standardValues.DryContactStatus),
		CutOffDeviceStatus:                 float64(standardValues.CutOffDeviceStatus),
		LinkyTerminalShieldStatus:          float64(standardValues.LinkyTerminalShieldStatus),
		SurgeStatus:                        float64(standardValues.SurgeStatus),
		ReferencePowerExceededStatus:       float64(standardValues.ReferencePowerExceededStatus),
		ConsumptionStatus:                  float64(standardValues.ConsumptionStatus),
		EnergyDirectionStatus:              float64(standardValues.EnergyDirectionStatus),
		ContractTypePriceStatus:            float64(standardValues.ContractTypePriceStatus),
		ContractTypePriceDistributorStatus: float64(standardValues.ContractTypePriceDistributorStatus),
		ClockStatus:                        float64(standardValues.ClockStatus),
		TicStatus:                          float64(standardValues.TicStatus),
		EuridisLinkStatus:                  float64(standardValues.EuridisLinkStatus),
		CPLStatus:                          float64(standardValues.CPLStatus),
		CPLSyncStatus:                      float64(standardValues.CPLSyncStatus),
		TempoContractColorStatus:           float64(standardValues.TempoContractColorStatus),
		TempoContractNextDayColorStatus:    float64(standardValues.TempoContractNextDayColorStatus),
		MovingPeakNoticeStatus:             float64(standardValues.MovingPeakNoticeStatus),
		MovingPeakStatus:                   float64(standardValues.MovingPeakStatus),
		MovingPeakStart1:                   float64(standardValues.Dpm1),
		MovingPeakEnd1:                     float64(standardValues.Fpm1),
		MovingPeakStart2:                   float64(standardValues.Dpm2),
		MovingPeakEnd2:                     float64(standardValues.Fpm2),
		MovingPeakStart3:                   float64(standardValues.Dpm3),
		MovingPeakEnd3:                     float64(standardValues.Fpm3),
		Prm:                                standardValues.Prm,
		Relay1:                             float64(standardValues.Relai1),
		Relay2:                             float64(standardValues.Relai2),
		Relay3:                             float64(standardValues.Relai3),
		Relay4:                             float64(standardValues.Relai4),
		Relay5:                             float64(standardValues.Relai5),
		Relay6:                             float64(standardValues.Relai6),
		Relay7:                             float64(standardValues.Relai7),
		Relay8:                             float64(standardValues.Relai8),
		CurrentPricingNumber:               strconv.FormatInt(int64(standardValues.Ntarf), 10),
		ContractTypeDayNumber:              strconv.FormatInt(int64(standardValues.Njourf), 10),
		ContractTypeNextDayNumber:          strconv.FormatInt(int64(standardValues.Njourfnd), 10),
		ContractTypeNextDayProfile:         standardValues.Pjourfnd,
		PeakNextDayProfile:                 standardValues.Ppointe,
	}
}
