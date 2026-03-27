package scraper

import (
	"crypto/x509"
	"os"
	"testing"

	"github.com/rtoma/homewizard-prometheus-exporter/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewHTTPClientHTTPSWithoutCACert(t *testing.T) {
	device := &config.Device{
		Name:        "battery",
		Host:        "192.168.1.100",
		Port:        443,
		BearerToken: "token123",
	}

	client, err := NewHTTPClient(device, nil)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "https://192.168.1.100:443", client.baseURL)
	assert.Equal(t, "token123", client.bearerToken)
}

func TestNewHTTPClientHTTPWithoutCACert(t *testing.T) {
	device := &config.Device{
		Name: "socket",
		Host: "192.168.1.101",
		Port: 80,
	}

	client, err := NewHTTPClient(device, nil)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "http://192.168.1.101:80", client.baseURL)
}

func TestNewHTTPClientHTTPSWithValidCACert(t *testing.T) {
	device := &config.Device{
		Name:        "battery",
		Host:        "192.168.1.100",
		Port:        443,
		BearerToken: "token123",
	}

	// Load actual homewizard CA certificate
	caCertPEM, err := os.ReadFile("../../homewizard-ca.pem")
	assert.NoError(t, err)
	assert.NotEmpty(t, caCertPEM)

	// Create cert pool
	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCertPEM)
	assert.True(t, ok)

	client, err := NewHTTPClient(device, caCertPool)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "https://192.168.1.100:443", client.baseURL)
	assert.Equal(t, "token123", client.bearerToken)
}
