# Prometheus Watermeter exporter

Prometheus exporter for HomeWizard Wifi Watermeter.

## Installation

With Go install:

```shell
go install github.com/loafoe/prometheus-watermeter-exporter@latest
```

## Usage

```
Usage of ./prometheus-watermeter-exporter:
  -addr string
        IP address of Watermeter on your network
  -listen string
        Listen address for HTTP metrics (default "127.0.0.1:8880")
  -verbose
        Verbose output logging
     
```

### Example output

```
# HELP watermeter_active_liter_lpm Active liter usage per minute
# TYPE watermeter_active_liter_lpm gauge
watermeter_active_liter_lpm{serial="deadbeafabc"} 0
# HELP watermeter_total_liter_m3 Total liters in cubic meter
# TYPE watermeter_total_liter_m3 gauge
watermeter_total_liter_m3{serial="deadbeafabc"} 176.541
# HELP watermeter_total_liter_offset_m3 Total liter offset in cubic meter
# TYPE watermeter_total_liter_offset_m3 gauge
watermeter_total_liter_offset_m3{serial="deadbeafabc"} 0
# HELP watermeter_wifi_strength wifi signal strength
# TYPE watermeter_wifi_strength gauge
watermeter_wifi_strength{serial="deadbeafabc"} 100
```

## License

License is [MIT](LICENSE.md)