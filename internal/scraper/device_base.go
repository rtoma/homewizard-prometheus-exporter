package scraper

import (
	"log"
	"strings"

	"github.com/rtoma/homewizard-prometheus-exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	apiVersionHeader = "X-Api-Version"
	apiVersion2      = "2"
)

var (
	deviceWifiStrengthGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_device_wifi_strength",
			Help: "WiFi signal strength (percentage, 0-100)",
		},
		[]string{"name"},
	)
	deviceInfo = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_device_info",
			Help: "Device information",
		},
		[]string{"name", "product_type", "product_name", "serial", "firmware_version"},
	)
	// API v2 specific metrics
	deviceWifiRSSIDBGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_device_wifi_rssi_db",
			Help: "WiFi signal strength in dBm",
		},
		[]string{"name"},
	)
	deviceUptimeSecondsGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_device_uptime_seconds",
			Help: "Device uptime in seconds",
		},
		[]string{"name"},
	)
	deviceCloudEnabledGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_device_cloud_enabled",
			Help: "Cloud communication status (0=disabled, 1=enabled)",
		},
		[]string{"name"},
	)
	deviceStatusLedBrightnessGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_device_status_led_brightness_pct",
			Help: "Status LED brightness in percent",
		},
		[]string{"name"},
	)
)

type deviceBase struct {
	debugEnabled bool
	*config.Device
	*DeviceInfo
	httpClient *HTTPClient
	l          *log.Logger
}

// isAPIv2 returns true if device uses API v2
func (d *deviceBase) isAPIv2() bool {
	return strings.HasPrefix(d.APIVersion, "2.")
}

// getDataEndpoint returns the appropriate data endpoint for the API version
func (d *deviceBase) getDataEndpoint() string {
	if d.isAPIv2() {
		return "/api/measurement"
	}
	return "/api/v1/data"
}

// getV2Headers returns HTTP headers for API v2 requests
func (d *deviceBase) getV2Headers() map[string]string {
	return map[string]string{apiVersionHeader: apiVersion2}
}

func (d *deviceBase) getData(target interface{}) error {
	headers := make(map[string]string)
	opName := "get-data" // v1
	if d.isAPIv2() {
		headers = d.getV2Headers()
		opName = "get-measurement"
	}

	if err := d.httpClient.Get(d.getDataEndpoint(), headers, target); err != nil {
		deviceScrapeErrors.WithLabelValues(d.Name, opName).Inc()
		return err
	}
	return nil
}

// getSystemData fetches system data from /api/system (API v2)
func (d *deviceBase) getSystemData(target *SystemData) error {
	if err := d.httpClient.Get("/api/system", d.getV2Headers(), target); err != nil {
		deviceScrapeErrors.WithLabelValues(d.Name, "get-system").Inc()
		return err
	}
	return nil
}

// setDeviceInfo sets the device info metric (always set to 1 for label identification)
func (d *deviceBase) setDeviceInfo() {
	deviceInfo.WithLabelValues(d.Name, d.ProductType, d.ProductName, d.Serial, d.FirmwareVersion).Set(1)
}

// collectGenericMetrics collects WiFi metrics from GenericData (API v1)
func (d *deviceBase) collectGenericMetrics(data GenericData) {
	deviceWifiStrengthGauge.WithLabelValues(d.Name).Set(data.WifiStrength)
	d.setDeviceInfo()
}

// collectSystemMetrics collects WiFi and system metrics from SystemData (API v2)
func (d *deviceBase) collectSystemMetrics(data SystemData) {
	// Convert dBm to percentage for backward compatibility
	wifiStrengthPct := convertDBmToPercentage(data.WifiRSSIDB)
	deviceWifiStrengthGauge.WithLabelValues(d.Name).Set(wifiStrengthPct)

	// Also expose raw dBm value
	deviceWifiRSSIDBGauge.WithLabelValues(d.Name).Set(data.WifiRSSIDB)

	// New v2-specific metrics
	deviceUptimeSecondsGauge.WithLabelValues(d.Name).Set(data.UptimeSeconds)
	deviceCloudEnabledGauge.WithLabelValues(d.Name).Set(boolToFloat(data.CloudEnabled))
	deviceStatusLedBrightnessGauge.WithLabelValues(d.Name).Set(data.StatusLedBrightnessPct)

	d.setDeviceInfo()
}

// collectDeviceInfo collects WiFi and device information metrics
// For v2 devices: fetches from /api/system and exposes additional metrics
// For v1 devices: uses data from GenericData embedded in measurement response
func (d *deviceBase) collectDeviceInfo(genericData *GenericData) {
	if d.isAPIv2() {
		systemData := &SystemData{}
		if err := d.getSystemData(systemData); err != nil {
			d.l.Printf("failed to fetch system data: %v", err)
			// Still set basic device info even if system data fetch fails
			d.setDeviceInfo()
		} else {
			if d.debugEnabled {
				d.l.Printf("DEBUG: System data: %+v", systemData)
			}
			d.collectSystemMetrics(*systemData)
		}
	} else {
		// v1 API: use GenericData from measurement response
		if genericData != nil {
			d.collectGenericMetrics(*genericData)
		} else {
			// Fallback if no generic data available
			d.setDeviceInfo()
		}
	}
}
