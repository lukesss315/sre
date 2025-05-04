package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 1. Counter
var requestCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "demo_request_total",
		Help: "Total number of demo requests",
	},
)

// 2. Gauge
var temperatureGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "demo_temperature_celsius",
		Help: "Current temperature in Celsius",
	},
)

// 3. Histogram
var requestDurationHistogram = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name:    "demo_request_duration_seconds",
		Help:    "Histogram of response time for demo requests",
		Buckets: prometheus.DefBuckets, // 默认桶：0.005~10s
	},
)

// 4. Summary
var requestDurationSummary = prometheus.NewSummary(
	prometheus.SummaryOpts{
		Name:       "demo_request_duration_summary_seconds",
		Help:       "Summary of response time for demo requests",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}, // P50/P90/P99
	},
)

func init() {
	// 注册指标
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(temperatureGauge)
	prometheus.MustRegister(requestDurationHistogram)
	prometheus.MustRegister(requestDurationSummary)
}

func main() {
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		for {
			// Counter 示例
			requestCounter.Inc()

			// Gauge 示例
			temperature := 20 + rand.Float64()*15 // 20~35°C
			temperatureGauge.Set(temperature)

			// 模拟请求耗时
			duration := rand.Float64() // 0~1秒
			requestDurationHistogram.Observe(duration)
			requestDurationSummary.Observe(duration)

			time.Sleep(2 * time.Second)
		}
	}()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
