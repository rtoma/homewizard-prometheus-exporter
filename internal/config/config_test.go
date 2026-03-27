package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEnvironmentMissingDevices(t *testing.T) {
	c := &Config{}
	err := parseEnvironment(c, []string{})
	assert.Contains(t, err.Error(), "no devices configured")
}

func TestParseEnvironmentDefault(t *testing.T) {
	c := &Config{}
	env := []string{
		"PROMETHEUS_LISTEN_ADDR=FOO",
		"DEVICES=name1:host1,name2:host2:8080",
		"INTERVAL=30",
	}
	err := parseEnvironment(c, env)
	assert.Empty(t, err)
	assert.Equal(t, "FOO", c.PrometheusListenAddr)
	assert.Equal(t, 2, len(c.Devices))
	assert.Equal(t, "name1", c.Devices[0].Name)
	assert.Equal(t, "host1", c.Devices[0].Host)
	assert.Equal(t, 80, c.Devices[0].Port)
	assert.Equal(t, "name2", c.Devices[1].Name)
	assert.Equal(t, "host2", c.Devices[1].Host)
	assert.Equal(t, 8080, c.Devices[1].Port)
	assert.Equal(t, 30, c.Interval)
}

func TestParseEnvironmentInvalidInterval(t *testing.T) {
	tests := []struct {
		name   string
		env    []string
		errmsg string
	}{
		{
			name: "invalid interval",
			env: []string{
				"DEVICES=name1:host1",
				"INTERVAL=foo",
			},
			errmsg: "expected a number, but its not",
		},
		{
			name: "interval too low",
			env: []string{
				"DEVICES=name1:host1",
				"INTERVAL=9",
			},
			errmsg: "expected a number >= 10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{}
			err := parseEnvironment(c, tt.env)
			assert.Contains(t, err.Error(), tt.errmsg)
		})
	}
}

func TestParseEnvironmentHTTPSDeviceWithToken(t *testing.T) {
	c := &Config{}
	env := []string{
		"DEVICES=battery:192.168.1.100:443:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
	}
	err := parseEnvironment(c, env)
	assert.Empty(t, err)
	assert.Equal(t, 1, len(c.Devices))
	assert.Equal(t, "battery", c.Devices[0].Name)
	assert.Equal(t, "192.168.1.100", c.Devices[0].Host)
	assert.Equal(t, 443, c.Devices[0].Port)
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", c.Devices[0].BearerToken)
	assert.True(t, c.Devices[0].IsHTTPS())
}

func TestParseEnvironmentHTTPSDeviceMissingToken(t *testing.T) {
	c := &Config{}
	env := []string{
		"DEVICES=battery:192.168.1.100:443",
	}
	err := parseEnvironment(c, env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "HTTPS devices (port 443) require bearer token")
}

func TestParseEnvironmentHTTPDeviceWithToken(t *testing.T) {
	c := &Config{}
	env := []string{
		"DEVICES=socket:192.168.1.101:80:sometoken",
	}
	err := parseEnvironment(c, env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bearer token only supported for HTTPS devices (port 443)")
}

func TestParseEnvironmentMixedDevices(t *testing.T) {
	c := &Config{}
	env := []string{
		"DEVICES=battery:192.168.1.100:443:token123,socket:192.168.1.101,p1:192.168.1.102:80",
	}
	err := parseEnvironment(c, env)
	assert.Empty(t, err)
	assert.Equal(t, 3, len(c.Devices))

	// Battery device (HTTPS with token)
	assert.Equal(t, "battery", c.Devices[0].Name)
	assert.Equal(t, "192.168.1.100", c.Devices[0].Host)
	assert.Equal(t, 443, c.Devices[0].Port)
	assert.Equal(t, "token123", c.Devices[0].BearerToken)
	assert.True(t, c.Devices[0].IsHTTPS())

	// Socket device (HTTP, default port)
	assert.Equal(t, "socket", c.Devices[1].Name)
	assert.Equal(t, "192.168.1.101", c.Devices[1].Host)
	assert.Equal(t, 80, c.Devices[1].Port)
	assert.Equal(t, "", c.Devices[1].BearerToken)
	assert.False(t, c.Devices[1].IsHTTPS())

	// P1 device (HTTP, explicit port)
	assert.Equal(t, "p1", c.Devices[2].Name)
	assert.Equal(t, "192.168.1.102", c.Devices[2].Host)
	assert.Equal(t, 80, c.Devices[2].Port)
	assert.Equal(t, "", c.Devices[2].BearerToken)
	assert.False(t, c.Devices[2].IsHTTPS())
}

func TestParseEnvironmentCACertFile(t *testing.T) {
	c := &Config{}
	env := []string{
		"DEVICES=battery:192.168.1.100:443:token123",
		"CA_CERT_FILE=/etc/ssl/certs/homewizard-ca.pem",
	}
	err := parseEnvironment(c, env)
	assert.Empty(t, err)
	assert.Equal(t, "/etc/ssl/certs/homewizard-ca.pem", c.CACertFile)
}

func TestParseEnvironmentInvalidDeviceFormat(t *testing.T) {
	tests := []struct {
		name   string
		env    []string
		errmsg string
	}{
		{
			name: "missing host",
			env: []string{
				"DEVICES=onlyname",
			},
			errmsg: "expected at least <name>:<host>",
		},
		{
			name: "invalid port",
			env: []string{
				"DEVICES=name:host:notaport",
			},
			errmsg: "invalid port",
		},
		{
			name: "port 443 without token",
			env: []string{
				"DEVICES=battery:192.168.1.100:443",
			},
			errmsg: "HTTPS devices (port 443) require bearer token",
		},
		{
			name: "token on non-443 port",
			env: []string{
				"DEVICES=device:192.168.1.100:8080:token",
			},
			errmsg: "bearer token only supported for HTTPS devices (port 443)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{}
			err := parseEnvironment(c, tt.env)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errmsg)
		})
	}
}

func TestDeviceIsHTTPS(t *testing.T) {
	tests := []struct {
		name     string
		port     int
		expected bool
	}{
		{"port 443 is HTTPS", 443, true},
		{"port 80 is not HTTPS", 80, false},
		{"port 8080 is not HTTPS", 8080, false},
		{"port 8443 is not HTTPS", 8443, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Device{Port: tt.port}
			assert.Equal(t, tt.expected, d.IsHTTPS())
		})
	}
}

func TestParseEnvironmentWithDevicePrefix(t *testing.T) {
	c := &Config{}
	env := []string{
		"DEVICE_1=socket:192.168.1.101",
		"DEVICE_2=p1:192.168.1.102:80",
		"DEVICE_BATTERY=battery:192.168.1.100:443:token123",
	}
	err := parseEnvironment(c, env)
	assert.Empty(t, err)
	assert.Equal(t, 3, len(c.Devices))

	// Note: order might vary due to map iteration
	names := make(map[string]bool)
	for _, d := range c.Devices {
		names[d.Name] = true
	}
	assert.True(t, names["socket"])
	assert.True(t, names["p1"])
	assert.True(t, names["battery"])
}

func TestParseEnvironmentMixedDevicesAndDevicePrefix(t *testing.T) {
	c := &Config{}
	env := []string{
		"DEVICES=socket:192.168.1.101,p1:192.168.1.102:80",
		"DEVICE_BATTERY=battery:192.168.1.100:443:token123",
		"DEVICE_WATER=water:192.168.1.103",
	}
	err := parseEnvironment(c, env)
	assert.Empty(t, err)
	assert.Equal(t, 4, len(c.Devices))

	// Check all device names are present
	names := make(map[string]bool)
	for _, d := range c.Devices {
		names[d.Name] = true
	}
	assert.True(t, names["socket"])
	assert.True(t, names["p1"])
	assert.True(t, names["battery"])
	assert.True(t, names["water"])
}

func TestParseEnvironmentDevicePrefixInvalid(t *testing.T) {
	c := &Config{}
	env := []string{
		"DEVICE_1=invalid",
	}
	err := parseEnvironment(c, env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid 'DEVICE_1' env var value")
	assert.Contains(t, err.Error(), "expected at least <name>:<host>")
}

func TestParseEnvironmentDevicePrefixEmpty(t *testing.T) {
	c := &Config{}
	env := []string{
		"DEVICE_1=  ",
		"DEVICE_2=valid:192.168.1.100",
	}
	err := parseEnvironment(c, env)
	assert.Empty(t, err)
	assert.Equal(t, 1, len(c.Devices))
	assert.Equal(t, "valid", c.Devices[0].Name)
}
