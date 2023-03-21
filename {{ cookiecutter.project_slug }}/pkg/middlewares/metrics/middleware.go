package metrics

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const serviceName = "cp-service"

var (
	labels = prometheus.Labels{"service": serviceName}

	KafkaConsumeSuccessCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:        prometheus.BuildFQName("cp", "kafka", "consume_success"),
			Help:        "记录kafka消息消费成功的次数",
			ConstLabels: labels,
		}, []string{"topic", "bid"},
	)
	KafkaConsumeErrorCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:        prometheus.BuildFQName("cp", "kafka", "consume_error"),
			Help:        "记录kafka消息消费失败的次数",
			ConstLabels: labels,
		}, []string{"topic", "bid"},
	)
	KafkaConsumeProcessingCounter = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        prometheus.BuildFQName("cp", "kafka", "consume_processing"),
		Help:        "当前正在处理的kafka消息个数",
		ConstLabels: labels,
	}, []string{"topic", "bid"},
	)

	requestInProgressTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        prometheus.BuildFQName("cp", "http", "requests_in_progress_total"),
		Help:        "All the requests in progress",
		ConstLabels: labels,
	}, []string{"method"})

	requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        prometheus.BuildFQName("cp", "http", "requests_total"),
		Help:        "Count all http requests by status code, method and path.",
		ConstLabels: labels,
	},
		[]string{"status_code", "method", "path"},
	)
	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        prometheus.BuildFQName("cp", "http", "request_duration_seconds"),
		Help:        "Duration of all HTTP requests by status code, method and path.",
		ConstLabels: labels,
		Buckets: []float64{
			0.000000001, // 1ns
			0.000000002,
			0.000000005,
			0.00000001, // 10ns
			0.00000002,
			0.00000005,
			0.0000001, // 100ns
			0.0000002,
			0.0000005,
			0.000001, // 1µs
			0.000002,
			0.000005,
			0.00001, // 10µs
			0.00002,
			0.00005,
			0.0001, // 100µs
			0.0002,
			0.0005,
			0.001, // 1ms
			0.002,
			0.005,
			0.01, // 10ms
			0.02,
			0.05,
			0.1, // 100 ms
			0.2,
			0.5,
			1.0, // 1s
			2.0,
			5.0,
			10.0, // 10s
			15.0,
			20.0,
			30.0,
		},
	},
		[]string{"status_code", "method", "path"},
	)
)

func NewMiddleware(cfgs ...Config) fiber.Handler {
	cfg := configDefault(cfgs...)
	return func(ctx *fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(ctx) {
			return ctx.Next()
		}
		start := time.Now()
		method := ctx.Route().Method

		if ctx.Route().Path == "/metrics" {
			return ctx.Next()
		}

		requestInProgressTotal.WithLabelValues(method).Inc()
		defer requestInProgressTotal.WithLabelValues(method).Dec()

		err := ctx.Next()
		// initialize with default error code
		// https://docs.gofiber.io/guide/error-handling
		status := fiber.StatusInternalServerError
		if err != nil {
			if e, ok := err.(*fiber.Error); ok {
				// Get correct error code from fiber.Error type
				status = e.Code
			}
		} else {
			status = ctx.Response().StatusCode()
		}

		path := ctx.Route().Path

		statusCode := strconv.Itoa(status)
		requestsTotal.WithLabelValues(statusCode, method, path).Inc()

		elapsed := float64(time.Since(start).Nanoseconds()) / 1e9
		requestDuration.WithLabelValues(statusCode, method, path).Observe(elapsed)

		return err
	}
}
