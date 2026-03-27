package main

import (
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rtoma/homewizard-prometheus-exporter/internal/config"
	"github.com/rtoma/homewizard-prometheus-exporter/internal/logger"
	"github.com/rtoma/homewizard-prometheus-exporter/internal/scraper"
	"github.com/go-co-op/gocron/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	l := logger.NewLogger("main")
	l.Print("Starting homewizard-prometheus-exporter")

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	if cfg.DebugEnabled {
		l.Print("Debug enabled")
	}

	// Check if any device uses HTTPS
	hasHTTPSDevice := false
	for _, device := range cfg.Devices {
		if device.IsHTTPS() {
			hasHTTPSDevice = true
			break
		}
	}

	// Load CA certificate and create cert pool if provided
	var caCertPool *x509.CertPool
	if cfg.CACertFile != "" {
		caCertPEM, err := os.ReadFile(cfg.CACertFile)
		if err != nil {
			log.Fatalf("Failed to read CA cert file %s: %v", cfg.CACertFile, err)
		}
		caCertPool = x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCertPEM) {
			log.Fatal("Failed to parse CA certificate")
		}
		l.Printf("Loaded CA certificate from %s", cfg.CACertFile)
	} else if hasHTTPSDevice {
		log.Fatal("HTTPS device(s) configured but no CA certificate provided (set CA_CERT_FILE environment variable)")
	}

	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal(err)
	}

	l.Print("Start scheduler")
	go s.Start()

	l.Printf("Scrape interval: %s", time.Duration(cfg.Interval)*time.Second)

	go func() {
		for _, device := range cfg.Devices {
			l.Printf("Add device: %q on %s:%d", device.Name, device.Host, device.Port)
			_, err = s.NewJob(
				gocron.DurationJob(time.Duration(cfg.Interval)*time.Second),
				gocron.NewTask(func() {
					scraper.Scrape(device, caCertPool, cfg.DebugEnabled)
				}),
				gocron.WithStartAt(gocron.WithStartImmediately()),
			)
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second) // splay a little
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, "ok")
	})
	log.Printf("Starting HTTP server on %s", cfg.PrometheusListenAddr)
	log.Fatal(http.ListenAndServe(cfg.PrometheusListenAddr, nil))
}
