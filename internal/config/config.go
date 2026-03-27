package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rtoma/homewizard-prometheus-exporter/internal/logger"
)

const DEFAULT_PORT = 80

type Config struct {
	PrometheusListenAddr string
	Interval             int
	Devices              []*Device
	DebugEnabled         bool
	CACertFile           string
}

type Device struct {
	// Device name
	Name string
	// Hostname or IP address of device
	Host string
	// Port number of device API HTTP listener
	Port int
	// Bearer token for authentication (HTTPS devices only)
	BearerToken string
}

// IsHTTPS returns true if device uses HTTPS (port 443)
func (d *Device) IsHTTPS() bool {
	return d.Port == 443
}

var log = logger.NewLogger("config")

// parseDeviceString parses a device string in format: name:host[:port][:bearer_token]
// Returns a Device or an error if the format is invalid
func parseDeviceString(deviceStr string) (*Device, error) {
	parts := strings.Split(deviceStr, ":")
	if len(parts) < 2 {
		return nil, fmt.Errorf("expected at least <name>:<host>, got: %s", deviceStr)
	}

	d := &Device{
		Name:        parts[0],
		Host:        parts[1],
		Port:        DEFAULT_PORT,
		BearerToken: "",
	}

	// Parse optional port
	if len(parts) >= 3 && parts[2] != "" {
		num, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid port in %s: %v", deviceStr, err)
		}
		d.Port = num
	}

	// Parse bearer token (only valid for port 443)
	if len(parts) >= 4 {
		if d.Port != 443 {
			return nil, fmt.Errorf("bearer token only supported for HTTPS devices (port 443): %s", deviceStr)
		}
		d.BearerToken = parts[3]
	}

	// Validate: port 443 requires bearer token
	if d.Port == 443 && d.BearerToken == "" {
		return nil, fmt.Errorf("HTTPS devices (port 443) require bearer token in format name:host:443:token, got: %s", deviceStr)
	}

	return d, nil
}

func NewConfig() (*Config, error) {
	c := &Config{
		PrometheusListenAddr: ":8080",
		Interval:             10,
		Devices:              []*Device{},
		DebugEnabled:         false,
	}

	if err := parseEnvironment(c, os.Environ()); err != nil {
		return nil, err
	}

	return c, nil
}

func parseEnvironment(c *Config, env []string) error {
	// transform env []string into dict
	// where each item has a key and value join by equal character
	envvars := make(map[string]string)
	for _, envvar := range env {
		parts := strings.SplitN(envvar, "=", 2)
		if len(parts) != 2 {
			log.Printf("WARN:gnoring invalid env var: %s", envvar)
			continue
		}
		envvars[parts[0]] = parts[1]
	}

	if val, ok := envvars["PROMETHEUS_LISTEN_ADDR"]; ok {
		c.PrometheusListenAddr = val
	}

	if val, ok := envvars["INTERVAL"]; ok {
		num, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("invalid 'INTERVAL' env var value: expected a number, but its not: %s: %v", val, err)
		}
		if num < 10 {
			return fmt.Errorf("invalid 'INTERVAL' env var value: expected a number >= 10, but its not: %d", num)
		}
		c.Interval = num
	}

	// Parse DEVICES env var (comma-separated list)
	if val, ok := envvars["DEVICES"]; ok {
		for _, deviceStr := range strings.Split(val, ",") {
			deviceStr = strings.TrimSpace(deviceStr)
			if deviceStr == "" {
				continue
			}
			d, err := parseDeviceString(deviceStr)
			if err != nil {
				return fmt.Errorf("invalid 'DEVICES' env var value: %v", err)
			}
			c.Devices = append(c.Devices, d)
		}
	}

	// Parse DEVICE_* env vars (individual device definitions)
	for key, val := range envvars {
		if strings.HasPrefix(key, "DEVICE_") {
			val = strings.TrimSpace(val)
			if val == "" {
				continue
			}
			d, err := parseDeviceString(val)
			if err != nil {
				return fmt.Errorf("invalid '%s' env var value: %v", key, err)
			}
			c.Devices = append(c.Devices, d)
		}
	}

	// Validate that at least one device is configured
	if len(c.Devices) == 0 {
		return fmt.Errorf("no devices configured: set DEVICES or DEVICE_* env vars")
	}

	if val, ok := envvars["DEBUG_ENABLED"]; ok && val != "" {
		c.DebugEnabled = true
	}

	if val, ok := envvars["CA_CERT_FILE"]; ok {
		c.CACertFile = val
	}

	return nil
}
