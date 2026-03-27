package scraper

import (
	"fmt"
	"math"

	"github.com/rtoma/homewizard-prometheus-exporter/internal/config"
	"github.com/rtoma/homewizard-prometheus-exporter/internal/logger"
)

const (
	// WiFi signal strength conversion constants
	wifiExcellentDBm = -30.0 // Signal strength considered 100%
	wifiPoorDBm      = -90.0 // Signal strength considered 0%
)

type GenericData struct {
	WifiSSID     string  `json:"wifi_ssid"`
	WifiStrength float64 `json:"wifi_strength"`
}

type SystemData struct {
	WifiSSID               string  `json:"wifi_ssid"`
	WifiRSSIDB             float64 `json:"wifi_rssi_db"`
	UptimeSeconds          float64 `json:"uptime_s"`
	CloudEnabled           bool    `json:"cloud_enabled"`
	StatusLedBrightnessPct float64 `json:"status_led_brightness_pct"`
}

// convertDBmToPercentage converts WiFi signal strength from dBm to percentage
// Using a linear scale where excellent signal (-30 dBm) = 100% and poor signal (-90 dBm) = 0%
func convertDBmToPercentage(dBm float64) float64 {
	if dBm >= wifiExcellentDBm {
		return 100
	}
	if dBm <= wifiPoorDBm {
		return 0
	}
	// Linear interpolation between excellent and poor thresholds
	dbmRange := wifiExcellentDBm - wifiPoorDBm
	return math.Round((dBm - wifiPoorDBm) * 100 / dbmRange)
}

// boolToFloat converts a boolean to float64 (0.0 or 1.0) for Prometheus gauges
func boolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

type DeviceScraper interface {
	Scrape()
}

func NewDeviceScraper(device *config.Device, info *DeviceInfo, httpClient *HTTPClient, debugEnabled bool) (DeviceScraper, error) {
	log := logger.NewLogger(fmt.Sprintf("scrape:%s/%s", info.ProductName, device.Name))

	switch info.ProductType {
	case "HWE-WTR":
		return &WatermeterScraper{deviceBase: deviceBase{debugEnabled, device, info, httpClient, log}}, nil
	case "HWE-P1":
		return &P1MeterScraper{deviceBase: deviceBase{debugEnabled, device, info, httpClient, log}}, nil
	case "HWE-SKT":
		return &EnergySocketScraper{deviceBase: deviceBase{debugEnabled, device, info, httpClient, log}}, nil
	case "HWE-BAT":
		return &BatteryScraper{deviceBase: deviceBase{debugEnabled, device, info, httpClient, log}}, nil
	}
	return nil, fmt.Errorf("unknown device, info: %+v", info)
}
