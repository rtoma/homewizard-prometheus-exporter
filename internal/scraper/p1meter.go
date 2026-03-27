package scraper

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	p1MeterActiveTariff = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_active_tariff",
			Help: "Active tariff",
		},
		[]string{"name"},
	)
	p1MeterTotalPowerImportKWh = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_total_power_import_kwh",
			Help: "The energy usage meter reading for all tariffs in kWh",
		},
		[]string{"name"},
	)
	p1MeterTotalPowerImportT1KWh = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_total_power_import_t1_kwh",
			Help: "The energy usage meter reading for tariff 1 in kWh",
		},
		[]string{"name"},
	)
	p1MeterTotalPowerImportT2KWh = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_total_power_import_t2_kwh",
			Help: "The energy usage meter reading for tariff 2 in kWh",
		},
		[]string{"name"},
	)
	p1MeterTotalPowerExportKWh = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_total_power_export_kwh",
			Help: "The energy feed-in meter reading for all tariffs in kWh",
		},
		[]string{"name"},
	)
	p1MeterTotalPowerExportT1KWh = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_total_power_export_t1_kwh",
			Help: "The energy feed-in meter reading for tariff 1 in kWh",
		},
		[]string{"name"},
	)
	p1MeterTotalPowerExportT2KWh = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_total_power_export_t2_kwh",
			Help: "The energy feed-in meter reading for tariff 2 in kWh",
		},
		[]string{"name"},
	)
	p1MeterActivePowerW = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_active_power_w",
			Help: "The total active usage in watt",
		},
		[]string{"name"},
	)
	p1MeterActivePowerL1W = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_active_power_l1_w",
			Help: "The active usage for phase 1 in watt",
		},
		[]string{"name"},
	)
	p1MeterActivePowerL2W = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_active_power_l2_w",
			Help: "The active usage for phase 2 in watt",
		},
		[]string{"name"},
	)
	p1MeterActivePowerL3W = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_active_power_l3_w",
			Help: "The active usage for phase 3 in watt",
		},
		[]string{"name"},
	)
	p1MeterActiveVoltageL1V = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_active_voltage_l1_v",
			Help: "The active voltage for phase 1 in volt",
		},
		[]string{"name"},
	)
	p1MeterActiveVoltageL2V = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_active_voltage_l2_v",
			Help: "The active voltage for phase 2 in volt",
		},
		[]string{"name"},
	)
	p1MeterActiveVoltageL3V = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_active_voltage_l3_v",
			Help: "The active voltage for phase 3 in volt",
		},
		[]string{"name"},
	)
	p1MeterActiveCurrentA = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_active_current_a",
			Help: "The active current in ampere",
		},
		[]string{"name"},
	)
	p1MeterActiveCurrentL1A = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_active_current_l1_a",
			Help: "The active current for phase 1 in ampere",
		},
		[]string{"name"},
	)
	p1MeterActiveCurrentL2A = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_active_current_l2_a",
			Help: "The active current for phase 2 in ampere",
		},
		[]string{"name"},
	)
	p1MeterActiveCurrentL3A = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_active_current_l3_a",
			Help: "The active current for phase 3 in ampere",
		},
		[]string{"name"},
	)
	p1MeterVoltageSagL1Count = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_voltage_sag_l1_count",
			Help: "Number of voltage sags detected by meter for phase 1",
		},
		[]string{"name"},
	)
	p1MeterVoltageSagL2Count = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_voltage_sag_l2_count",
			Help: "Number of voltage sags detected by meter for phase 2",
		},
		[]string{"name"},
	)
	p1MeterVoltageSagL3Count = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_voltage_sag_l3_count",
			Help: "Number of voltage sags detected by meter for phase 3",
		},
		[]string{"name"},
	)
	p1MeterVoltageSwellL1Count = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_voltage_swell_l1_count",
			Help: "Number of voltage swells detected by meter for phase 1",
		},
		[]string{"name"},
	)
	p1MeterVoltageSwellL2Count = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_voltage_swell_l2_count",
			Help: "Number of voltage swells detected by meter for phase 2",
		},
		[]string{"name"},
	)
	p1MeterVoltageSwellL3Count = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_voltage_swell_l3_count",
			Help: "Number of voltage swells detected by meter for phase 3",
		},
		[]string{"name"},
	)
	p1MeterAnyPowerFailCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_any_power_fail_count",
			Help: "Number of power failures detected by meter",
		},
		[]string{"name"},
	)
	p1MeterLongPowerFailCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_long_power_fail_count",
			Help: "Number of 'long' power fails detected by meter.",
		},
		[]string{"name"},
	)
	p1MeterDeviceInfo = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_p1meter_device_info",
			Help: "Information about the device",
		},
		[]string{"name", "smr_version", "meter_model", "unique_id"},
	)
)

type P1MeterData struct {
	GenericData
	SmrVersion            int     `json:"smr_version"`
	MeterModel            string  `json:"meter_model"`
	UniqueID              string  `json:"unique_id"`
	ActiveTariff          int     `json:"active_tariff"`
	TotalPowerImportKWh   float64 `json:"total_power_import_kwh"`
	TotalPowerImportT1KWh float64 `json:"total_power_import_t1_kwh"`
	TotalPowerImportT2KWh float64 `json:"total_power_import_t2_kwh"`
	TotalPowerExportKWh   float64 `json:"total_power_export_kwh"`
	TotalPowerExportT1KWh float64 `json:"total_power_export_t1_kwh"`
	TotalPowerExportT2KWh float64 `json:"total_power_export_t2_kwh"`
	ActivePowerW          float64 `json:"active_power_w"`
	ActivePowerL1W        float64 `json:"active_power_l1_w"`
	ActivePowerL2W        float64 `json:"active_power_l2_w"`
	ActivePowerL3W        float64 `json:"active_power_l3_w"`
	ActiveVoltageL1V      float64 `json:"active_voltage_l1_v"`
	ActiveVoltageL2V      float64 `json:"active_voltage_l2_v"`
	ActiveVoltageL3V      float64 `json:"active_voltage_l3_v"`
	ActiveCurrentA        float64 `json:"active_current_a"`
	ActiveCurrentL1A      float64 `json:"active_current_l1_a"`
	ActiveCurrentL2A      float64 `json:"active_current_l2_a"`
	ActiveCurrentL3A      float64 `json:"active_current_l3_a"`
	VoltageSagL1Count     float64 `json:"voltage_sag_l1_count"`
	VoltageSagL2Count     float64 `json:"voltage_sag_l2_count"`
	VoltageSagL3Count     float64 `json:"voltage_sag_l3_count"`
	VoltageSwellL1Count   float64 `json:"voltage_swell_l1_count"`
	VoltageSwellL2Count   float64 `json:"voltage_swell_l2_count"`
	VoltageSwellL3Count   float64 `json:"voltage_swell_l3_count"`
	AnyPowerFailCount     float64 `json:"any_power_fail_count"`
	LongPowerFailCount    float64 `json:"long_power_fail_count"`
	// External              []string `json:"external"`
}

type P1MeterScraper struct {
	deviceBase
}

func (s *P1MeterScraper) Scrape() {
	// s.l.Print("Scrape")
	data := &P1MeterData{}
	if err := s.getData(data); err != nil {
		s.l.Printf("failed to fetch measurement data: %v", err)
		return
	}
	if s.debugEnabled {
		s.l.Printf("DEBUG: Scraped data: %+v", data)
	}

	s.collectDeviceInfo(&data.GenericData)

	p1MeterActiveTariff.WithLabelValues(s.Name).Set(float64(data.ActiveTariff))
	p1MeterTotalPowerImportKWh.WithLabelValues(s.Name).Set(data.TotalPowerImportKWh)
	p1MeterTotalPowerImportT1KWh.WithLabelValues(s.Name).Set(data.TotalPowerImportT1KWh)
	p1MeterTotalPowerImportT2KWh.WithLabelValues(s.Name).Set(data.TotalPowerImportT2KWh)
	p1MeterTotalPowerExportKWh.WithLabelValues(s.Name).Set(data.TotalPowerExportKWh)
	p1MeterTotalPowerExportT1KWh.WithLabelValues(s.Name).Set(data.TotalPowerExportT1KWh)
	p1MeterTotalPowerExportT2KWh.WithLabelValues(s.Name).Set(data.TotalPowerExportT2KWh)
	p1MeterActivePowerW.WithLabelValues(s.Name).Set(data.ActivePowerW)
	p1MeterActivePowerL1W.WithLabelValues(s.Name).Set(data.ActivePowerL1W)
	p1MeterActivePowerL2W.WithLabelValues(s.Name).Set(data.ActivePowerL2W)
	p1MeterActivePowerL3W.WithLabelValues(s.Name).Set(data.ActivePowerL3W)
	p1MeterActiveVoltageL1V.WithLabelValues(s.Name).Set(data.ActiveVoltageL1V)
	p1MeterActiveVoltageL2V.WithLabelValues(s.Name).Set(data.ActiveVoltageL2V)
	p1MeterActiveVoltageL3V.WithLabelValues(s.Name).Set(data.ActiveVoltageL3V)
	p1MeterActiveCurrentA.WithLabelValues(s.Name).Set(data.ActiveCurrentA)
	p1MeterActiveCurrentL1A.WithLabelValues(s.Name).Set(data.ActiveCurrentL1A)
	p1MeterActiveCurrentL2A.WithLabelValues(s.Name).Set(data.ActiveCurrentL2A)
	p1MeterActiveCurrentL3A.WithLabelValues(s.Name).Set(data.ActiveCurrentL3A)
	p1MeterVoltageSagL1Count.WithLabelValues(s.Name).Set(data.VoltageSagL1Count)
	p1MeterVoltageSagL2Count.WithLabelValues(s.Name).Set(data.VoltageSagL2Count)
	p1MeterVoltageSagL3Count.WithLabelValues(s.Name).Set(data.VoltageSagL3Count)
	p1MeterVoltageSwellL1Count.WithLabelValues(s.Name).Set(data.VoltageSwellL1Count)
	p1MeterVoltageSwellL2Count.WithLabelValues(s.Name).Set(data.VoltageSwellL2Count)
	p1MeterVoltageSwellL3Count.WithLabelValues(s.Name).Set(data.VoltageSwellL3Count)
	p1MeterAnyPowerFailCount.WithLabelValues(s.Name).Set(data.AnyPowerFailCount)
	p1MeterLongPowerFailCount.WithLabelValues(s.Name).Set(data.LongPowerFailCount)
	p1MeterDeviceInfo.WithLabelValues(s.Name, strconv.Itoa(data.SmrVersion),
		data.MeterModel, data.UniqueID).Set(1)
}
