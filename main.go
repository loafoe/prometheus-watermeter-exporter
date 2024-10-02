package main

import (
	"flag"
	"github.com/loafoe/prometheus-watermeter-exporter/watermeter"
	"github.com/sirupsen/logrus"
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
	totalLiterM3 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricNamePrefix + "total_liter_m3",
		Help: "Total liters in cubic meter",
	})
	activeLiterPerMinute = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricNamePrefix + "active_liter_lpm",
		Help: "Active liter usage per minute",
	})
	totalLiterOffsetM3 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricNamePrefix + "total_liter_offset_m3",
		Help: "Total liter offset in cubic meter",
	})
	wifiStrength = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricNamePrefix + "wifi_strength",
		Help: "wifi signal strength",
	})
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
	flag.StringVar(&watermeterAddr, "ip", "", "IP address of Watermeter on your network")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output logging")
	flag.Parse()

	wm, err := watermeter.New(watermeterAddr)
	if err != nil {
		logger.Error("Quitting because of error opening watermeter address", "error", err, "addr", watermeterAddr)
		os.Exit(1)
	}

	// Start
	wm.Start()

	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	go func() {
		for t := range wm.Incoming {
			totalLiterM3.Set(t.Data.TotalLiterM3)
			activeLiterPerMinute.Set(float64(t.Data.ActiveLiterLpm))
			totalLiterOffsetM3.Set(float64(t.Data.TotalLiterOffsetM3))
			wifiStrength.Set(float64(t.Data.WifiStrength))
			time.Sleep(2 * time.Second)
		}
	}()

	logrus.Infoln("Start listening at", listenAddr)
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	logrus.Fatalln(http.ListenAndServe(listenAddr, nil))
}
