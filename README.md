# homewizard-prometheus-exporter

Prometheus exporter for [HomeWizard](https://www.homewizard.com/) devices. Scrapes the local device API and exposes metrics for Prometheus.

## Supported Devices

| Device | Product type | Protocol |
|---|---|---|
| P1 Meter (HWE-P1) | `HWE-P1` | HTTP (port 80) or HTTPS (port 443) |
| Energy Socket (HWE-SKT) | `HWE-SKT` | HTTP (port 80) or HTTPS (port 443) |
| Water Meter (HWE-WTR) | `HWE-WTR` | HTTP (port 80) |
| Battery (HWE-BAT) | `HWE-BAT` | HTTPS (port 443) |

The device type is detected automatically via the HomeWizard local API.

## Configuration

All configuration is done via environment variables.

| Variable | Default | Description |
|---|---|---|
| `DEVICES` | — | Comma-separated list of devices: `name:host[:port][:bearer_token]` |
| `DEVICE_*` | — | Individual device, e.g. `DEVICE_P1=p1:192.168.0.10` |
| `INTERVAL` | `10` | Scrape interval in seconds (minimum: 10) |
| `PROMETHEUS_LISTEN_ADDR` | `:8080` | Address the HTTP server listens on |
| `CA_CERT_FILE` | — | Path to CA certificate PEM file (required for HTTPS devices) |
| `DEBUG_ENABLED` | `false` | Enable debug logging |

At least one device must be configured via `DEVICES` or `DEVICE_*`.

### Device string format

```
name:host[:port][:bearer_token]
```

- **HTTP devices** (port 80, default): `p1:192.168.0.10`
- **HTTPS devices** (port 443): `battery:192.168.0.20:443:your-bearer-token`

### HTTPS setup

Newer HomeWizard devices use HTTPS with a bearer token. To connect:

1. Obtain the bearer token from the HomeWizard app (Settings → Devices → Local API).
2. Set `CA_CERT_FILE` to point to the HomeWizard CA certificate (see below).

**HomeWizard CA certificate**

The official HomeWizard Appliance Access CA certificate is included in this repository as [`homewizard-ca.pem`](homewizard-ca.pem). It is the public CA used to sign all HomeWizard device certificates and is safe to distribute. It expires **16 December 2031**.

```bash
export CA_CERT_FILE=homewizard-ca.pem
```

> **Note:** HomeWizard device TLS certificates use the device serial number as the Common Name, not the IP address or hostname. TLS hostname verification is therefore skipped; only the certificate chain is validated using the provided CA certificate.

## Usage

### Run directly

```bash
export DEVICES="p1:192.168.0.10,socket:192.168.0.11,battery:192.168.0.12:443:your-token"
export INTERVAL=30
export CA_CERT_FILE=/path/to/homewizard-ca.pem  # only needed for HTTPS devices
go run cmd/main.go
```

Metrics are available at `http://localhost:8080/metrics`, health check at `http://localhost:8080/health`.

### Docker

```bash
docker build -t homewizard-prometheus-exporter .

docker run \
  -e DEVICES="p1:192.168.0.10,socket:192.168.0.11" \
  -e INTERVAL=30 \
  -p 8080:8080 \
  homewizard-prometheus-exporter
```

For HTTPS devices, mount the CA certificate into the container:

```bash
docker run \
  -e DEVICES="battery:192.168.0.12:443:your-token" \
  -e CA_CERT_FILE=/certs/homewizard-ca.pem \
  -v /path/to/certs:/certs:ro \
  -p 8080:8080 \
  homewizard-prometheus-exporter
```

### Docker Compose example

```yaml
services:
  homewizard-exporter:
    build: .
    environment:
      DEVICES: "p1:192.168.0.10,socket:192.168.0.11"
      INTERVAL: "30"
    ports:
      - "8080:8080"
    restart: unless-stopped
```

## Metrics

All metrics carry a `name` label matching the device name you configured.

### Common (all devices)

| Metric | Description |
|---|---|
| `homewizard_device_info` | Device info (labels: product_type, product_name, serial, firmware_version) |
| `homewizard_device_wifi_strength` | WiFi signal strength (0–100%) |
| `homewizard_device_wifi_rssi_db` | WiFi RSSI in dBm (API v2 devices) |
| `homewizard_device_uptime_seconds` | Device uptime in seconds (API v2 devices) |
| `homewizard_device_cloud_enabled` | Cloud communication enabled (0/1, API v2 devices) |
| `homewizard_device_status_led_brightness_pct` | Status LED brightness % (API v2 devices) |
| `homewizard_device_scrape_latency_sec` | Duration of each scrape request |
| `homewizard_device_scrape_errors_total` | Total number of scrape errors |

### P1 Meter (HWE-P1)

| Metric | Description |
|---|---|
| `homewizard_p1meter_active_tariff` | Active tariff |
| `homewizard_p1meter_total_power_import_kwh` | Total energy usage (all tariffs) in kWh |
| `homewizard_p1meter_total_power_import_t1_kwh` | Energy usage tariff 1 in kWh |
| `homewizard_p1meter_total_power_import_t2_kwh` | Energy usage tariff 2 in kWh |
| `homewizard_p1meter_total_power_export_kwh` | Total energy feed-in (all tariffs) in kWh |
| `homewizard_p1meter_total_power_export_t1_kwh` | Energy feed-in tariff 1 in kWh |
| `homewizard_p1meter_total_power_export_t2_kwh` | Energy feed-in tariff 2 in kWh |
| `homewizard_p1meter_active_power_w` | Total active power in W |
| `homewizard_p1meter_active_power_l1_w` | Active power phase 1 in W |
| `homewizard_p1meter_active_power_l2_w` | Active power phase 2 in W |
| `homewizard_p1meter_active_power_l3_w` | Active power phase 3 in W |
| `homewizard_p1meter_active_voltage_l1_v` | Voltage phase 1 in V |
| `homewizard_p1meter_active_voltage_l2_v` | Voltage phase 2 in V |
| `homewizard_p1meter_active_voltage_l3_v` | Voltage phase 3 in V |
| `homewizard_p1meter_active_current_a` | Total active current in A |
| `homewizard_p1meter_active_current_l1_a` | Current phase 1 in A |
| `homewizard_p1meter_active_current_l2_a` | Current phase 2 in A |
| `homewizard_p1meter_active_current_l3_a` | Current phase 3 in A |

### Energy Socket (HWE-SKT)

| Metric | Description |
|---|---|
| `homewizard_energysocket_total_power_import_kwh` | Total energy usage in kWh |
| `homewizard_energysocket_total_power_export_kwh` | Total energy feed-in in kWh |
| `homewizard_energysocket_active_power_w` | Active power in W |
| `homewizard_energysocket_active_voltage_v` | Voltage in V |
| `homewizard_energysocket_active_current_a` | Current in A |
| `homewizard_energysocket_active_frequency_hz` | Frequency in Hz |

### Water Meter (HWE-WTR)

| Metric | Description |
|---|---|
| `homewizard_watermeter_active_liter_lpm` | Current flow in liters per minute |
| `homewizard_watermeter_total_liter_m3` | Total water usage in m³ |

### Battery (HWE-BAT)

| Metric | Description |
|---|---|
| `homewizard_battery_energy_import_kwh` | Total energy charged in kWh |
| `homewizard_battery_energy_export_kwh` | Total energy discharged in kWh |
| `homewizard_battery_power_w` | Active power in W (positive = charging, negative = discharging) |
| `homewizard_battery_voltage_v` | Battery voltage in V |
| `homewizard_battery_current_a` | Current in A |
| `homewizard_battery_frequency_hz` | Line frequency in Hz |
| `homewizard_battery_state_of_charge_pct` | State of charge in % |
| `homewizard_battery_cycles` | Total charge/discharge cycles |

## License

MIT — see [LICENSE](LICENSE).