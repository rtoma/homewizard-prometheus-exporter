package scraper

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	energySocketTotalPowerImportKWh = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_energysocket_total_power_import_kwh",
			Help: "The energy usage meter reading in kWh",
		},
		[]string{"name"},
	)
	energySocketTotalPowerExportKWh = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_energysocket_total_power_export_kwh",
			Help: "The energy feed-in meter reading in kWh",
		},
		[]string{"name"},
	)
	energySocketActivePowerW = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_energysocket_active_power_w",
			Help: "The total active usage in Watt",
		},
		[]string{"name"},
	)
	energySocketActiveVoltageV = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_energysocket_active_voltage_v",
			Help: "The active voltage in Volt",
		},
		[]string{"name"},
	)
	energySocketActiveCurrentA = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_energysocket_active_current_a",
			Help: "The active current in Ampere",
		},
		[]string{"name"},
	)
	energySocketActiveFrequencyHz = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_energysocket_active_frequency_hz",
			Help: "The active frequency in Hertz",
		},
		[]string{"name"},
	)
)

type EnergySocketData struct {
	GenericData
	TotalPowerImportKWh float64 `json:"total_power_import_kwh"`
	// TotalPowerImportT1KWh float64 `json:"total_power_import_t1_kwh"`
	TotalPowerExportKWh float64 `json:"total_power_export_kwh"`
	// TotalPowerExportT1KWh float64 `json:"total_power_export_t1_kwh"`
	ActivePowerW float64 `json:"active_power_w"`
	// ActivePowerL1W        float64 `json:"active_power_l1_w"`
	ActiveVoltageV    float64 `json:"active_voltage_v"`
	ActiveCurrentA    float64 `json:"active_current_a"`
	ActiveFrequencyHz float64 `json:"active_frequency_hz"`
}

type EnergySocketScraper struct {
	deviceBase
}

func (s *EnergySocketScraper) Scrape() {
	// s.l.Print("Scrape")
	data := &EnergySocketData{}
	if err := s.getData(data); err != nil {
		s.l.Printf("failed to fetch measurement data: %v", err)
		return
	}
	if s.debugEnabled {
		s.l.Printf("DEBUG: Scraped data: %+v", data)
	}

	s.collectDeviceInfo(&data.GenericData)

	energySocketTotalPowerImportKWh.WithLabelValues(s.Name).Set(data.TotalPowerImportKWh)
	energySocketTotalPowerExportKWh.WithLabelValues(s.Name).Set(data.TotalPowerExportKWh)
	energySocketActivePowerW.WithLabelValues(s.Name).Set(data.ActivePowerW)
	energySocketActiveVoltageV.WithLabelValues(s.Name).Set(data.ActiveVoltageV)
	energySocketActiveCurrentA.WithLabelValues(s.Name).Set(data.ActiveCurrentA)
	energySocketActiveFrequencyHz.WithLabelValues(s.Name).Set(data.ActiveFrequencyHz)
}
