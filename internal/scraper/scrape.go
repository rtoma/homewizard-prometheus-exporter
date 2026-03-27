package scraper

import (
	"crypto/x509"
	"time"

	"github.com/rtoma/homewizard-prometheus-exporter/internal/config"
	"github.com/rtoma/homewizard-prometheus-exporter/internal/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type DeviceInfo struct {
	ProductName     string `json:"product_name"`
	ProductType     string `json:"product_type"`
	Serial          string `json:"serial"`
	FirmwareVersion string `json:"firmware_version"`
	APIVersion      string `json:"api_version"`
}

var (
	deviceScrapeLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "homewizard_device_scrape_latency_sec",
			Help: "Duration of scrape request in seconds",
		},
		[]string{"device"},
	)
	deviceScrapeErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "homewizard_device_scrape_errors_total",
			Help: "Total number of scrape errors",
		},
		[]string{"device", "error"},
	)
)

func Scrape(device *config.Device, caCertPool *x509.CertPool, debugEnabled bool) {
	log := logger.NewLogger("scrape:" + device.Name)

	start := time.Now()
	defer func() {
		deviceScrapeLatency.WithLabelValues(device.Name).Observe(
			float64(time.Since(start).Milliseconds()) / 1000)
	}()

	// Create device-specific HTTP client
	httpClient, err := NewHTTPClient(device, caCertPool)
	if err != nil {
		log.Printf("error creating HTTP client: %v", err)
		deviceScrapeErrors.WithLabelValues(device.Name, "create-http-client").Inc()
		return
	}

	// Get device info from /api endpoint
	info := &DeviceInfo{}
	if err := httpClient.Get("/api", nil, info); err != nil {
		log.Printf("error requesting device info: %v", err)
		deviceScrapeErrors.WithLabelValues(device.Name, "get-device-info").Inc()
		return
	}
	if debugEnabled {
		log.Printf("DEBUG: Device info: %+v", info)
	}

	scraper, err := NewDeviceScraper(device, info, httpClient, debugEnabled)
	if err != nil {
		log.Printf("error creating device scraper: %v", err)
		deviceScrapeErrors.WithLabelValues(device.Name, "create-device-scraper").Inc()
		return
	}

	scraper.Scrape()
}
