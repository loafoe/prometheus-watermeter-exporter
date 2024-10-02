package main

import (
	"flag"
	"github.com/loafoe/prometheus-watermeter-exporter/watermeter"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
)

var listenAddr string
var watermeterAddr string
var verbose bool
var metricNamePrefix = "watermeter_"

var (
	registry     = prometheus.NewRegistry()
	totalLiterM3 = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricNamePrefix + "total_liter_m3",
		Help: "Total liters in cubic meter",
	}, []string{"serial"})
	activeLiterPerMinute = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricNamePrefix + "active_liter_lpm",
		Help: "Active liter usage per minute",
	}, []string{"serial"})
	totalLiterOffsetM3 = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricNamePrefix + "total_liter_offset_m3",
		Help: "Total liter offset in cubic meter",
	}, []string{"serial"})
	wifiStrength = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricNamePrefix + "wifi_strength",
		Help: "wifi signal strength",
	}, []string{"serial"})
)

func init() {
	registry.MustRegister(totalLiterM3)
	registry.MustRegister(activeLiterPerMinute)
	registry.MustRegister(totalLiterOffsetM3)
	registry.MustRegister(wifiStrength)
}

func main() {
	logger := slog.Default()
	flag.StringVar(&listenAddr, "listen", "127.0.0.1:8880", "Listen address for HTTP metrics")
	flag.StringVar(&watermeterAddr, "addr", "", "IP address of Watermeter on your network")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output logging")
	flag.Parse()

	logger.Info("Starting watermeter exporter", "address", watermeterAddr)

	wm, err := watermeter.New(watermeterAddr, verbose, logger)
	if err != nil {
		logger.Error("Quitting because of error opening watermeter address", "error", err, "addr", watermeterAddr)
		os.Exit(1)
	}

	// Start
	wm.Start()

	go func() {
		if verbose {
			logger.Info("Starting metrics updater")
		}
		for t := range wm.Incoming {
			if verbose {
				logger.Info("Received telegram", "info", t.Info, "data", t.Data)
			}
			totalLiterM3.WithLabelValues(t.Info.Serial).Set(t.Data.TotalLiterM3)
			activeLiterPerMinute.WithLabelValues(t.Info.Serial).Set(float64(t.Data.ActiveLiterLpm))
			totalLiterOffsetM3.WithLabelValues(t.Info.Serial).Set(float64(t.Data.TotalLiterOffsetM3))
			wifiStrength.WithLabelValues(t.Info.Serial).Set(float64(t.Data.WifiStrength))
			time.Sleep(2 * time.Second)
		}
		logger.Info("Metrics updater stopped")
		os.Exit(2)
	}()

	logger.Info("Start listening", "address", listenAddr)
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	_ = http.ListenAndServe(listenAddr, nil)
}
