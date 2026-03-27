package scraper

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rtoma/homewizard-prometheus-exporter/internal/config"
)

// HTTPClient handles HTTP/HTTPS requests with optional TLS and bearer authentication
type HTTPClient struct {
	client      *http.Client
	baseURL     string
	bearerToken string
}

// NewHTTPClient creates an HTTP client for a device with optional TLS and authentication
func NewHTTPClient(device *config.Device, caCertPool *x509.CertPool) (*HTTPClient, error) {
	var transport *http.Transport
	var protocol string

	if device.IsHTTPS() {
		protocol = "https"

		tlsConfig := &tls.Config{
			InsecureSkipVerify: true, // device common name is no IP/hostname
		}

		// Use provided cert pool if available
		if caCertPool != nil {
			tlsConfig.RootCAs = caCertPool
		}

		transport = &http.Transport{
			TLSClientConfig:   tlsConfig,
			DisableKeepAlives: true, // HomeWizard devices don't handle connection reuse
		}
	} else {
		protocol = "http"
		transport = &http.Transport{
			DisableKeepAlives: true,
		}
	}

	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}

	return &HTTPClient{
		client:      client,
		baseURL:     fmt.Sprintf("%s://%s:%d", protocol, device.Host, device.Port),
		bearerToken: device.BearerToken,
	}, nil
}

// Get performs an HTTP GET request with optional custom headers and bearer authentication
func (c *HTTPClient) Get(path string, headers map[string]string, target interface{}) error {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return err
	}

	// Add custom headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Add bearer token if configured
	if c.bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
