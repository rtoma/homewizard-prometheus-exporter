package scraper

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	batteryEnergyImportKWh = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_battery_energy_import_kwh",
			Help: "The energy usage meter reading in kWh",
		},
		[]string{"name"},
	)
	batteryEnergyExportKWh = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_battery_energy_export_kwh",
			Help: "The energy feed-in meter reading in kWh",
		},
		[]string{"name"},
	)
	batteryPowerW = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_battery_power_w",
			Help: "The total active usage in watt (positive = charging, negative = discharging)",
		},
		[]string{"name"},
	)
	batteryVoltageV = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_battery_voltage_v",
			Help: "The active voltage in volt",
		},
		[]string{"name"},
	)
	batteryCurrentA = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_battery_current_a",
			Help: "The active current in ampere",
		},
		[]string{"name"},
	)
	batteryFrequencyHz = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_battery_frequency_hz",
			Help: "Line frequency in hertz",
		},
		[]string{"name"},
	)
	batteryStateOfChargePct = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_battery_state_of_charge_pct",
			Help: "The current state of charge in percent",
		},
		[]string{"name"},
	)
	batteryCycles = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_battery_cycles",
			Help: "Total number of battery charge/discharge cycles",
		},
		[]string{"name"},
	)
)

// BatteryMeasurementData represents the data from API v2 /api/measurement endpoint
type BatteryMeasurementData struct {
	EnergyImportKWh  float64 `json:"energy_import_kwh"`
	EnergyExportKWh  float64 `json:"energy_export_kwh"`
	PowerW           float64 `json:"power_w"`
	VoltageV         float64 `json:"voltage_v"`
	CurrentA         float64 `json:"current_a"`
	FrequencyHz      float64 `json:"frequency_hz"`
	StateOfChargePct float64 `json:"state_of_charge_pct"`
	Cycles           float64 `json:"cycles"`
}

type BatteryScraper struct {
	deviceBase
}

func (s *BatteryScraper) Scrape() {
	data := &BatteryMeasurementData{}
	if err := s.getData(data); err != nil {
		s.l.Printf("failed to fetch measurement data: %v", err)
		return
	}
	if s.debugEnabled {
		s.l.Printf("DEBUG: Scraped data: %+v", data)
	}

	s.collectDeviceInfo(nil)

	// Collect battery-specific metrics
	batteryEnergyImportKWh.WithLabelValues(s.Name).Set(data.EnergyImportKWh)
	batteryEnergyExportKWh.WithLabelValues(s.Name).Set(data.EnergyExportKWh)
	batteryPowerW.WithLabelValues(s.Name).Set(data.PowerW)
	batteryVoltageV.WithLabelValues(s.Name).Set(data.VoltageV)
	batteryCurrentA.WithLabelValues(s.Name).Set(data.CurrentA)
	batteryFrequencyHz.WithLabelValues(s.Name).Set(data.FrequencyHz)
	batteryStateOfChargePct.WithLabelValues(s.Name).Set(data.StateOfChargePct)
	batteryCycles.WithLabelValues(s.Name).Set(data.Cycles)
}
