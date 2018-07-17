package tools

import (
	"time"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	historyBuckets = [...]float64{
		10., 20., 30., 50., 80., 100., 200., 300., 500., 1000., 2000., 3000.}
	DefaultMetricPath = "/helloworld/metrics"

	ResponseCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "helloworld_requests_total",
		Help: "Total request counts"}, []string{"method", "endpoint"})
	ErrorCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "helloworld_error_total",
		Help: "Total Error counts"}, []string{"method", "endpoint"})
	ResponseLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "helloworld_response_latency_millisecond",
		Help:    "Response latency (millisecond)",
		Buckets: historyBuckets[:]}, []string{"method", "endpoint"})
)

func init() {

	fmt.Println("prometheus_init .... ")
	prometheus.MustRegister(ResponseCounter)
	prometheus.MustRegister(ErrorCounter)
	prometheus.MustRegister(ResponseLatency)
}

func Prometheus() gin.HandlerFunc {
	return func(c *gin.Context) {
		endPoint := c.Request.URL.Path
		if endPoint == DefaultMetricPath {
			c.Next()
		} else {
			start := time.Now()
			method := c.Request.Method

			c.Next()

			elapsed := float64(time.Since(start).Nanoseconds()) / 1000000
			ResponseCounter.WithLabelValues(method, endPoint).Inc()
			ResponseLatency.WithLabelValues(method, endPoint).Observe(elapsed)
		}
	}
}

func LatestMetrics(c *gin.Context) {
	handler := promhttp.Handler()
	handler.ServeHTTP(c.Writer, c.Request)
}
