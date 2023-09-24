package metricecho

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// HTTPRecorderConfig lists all available configuration options for HTTPRecorder.
type HTTPRecorderConfig struct {
	// Namespace is the metrics namespace that will be added to all metric configurations. It will be a prefix to each
	// metric name. For example, if Namespace is "myapp", then requests_total metric will be myapp_http_requests_total
	// (after namespace there is also the subsystem prefix, "http" in this case).
	Namespace string

	// EnableTotalMetric whether to enable a metric to count the total number of http requests.
	EnableTotalMetric bool

	// EnableDurMetric whether to enable a metric to track the duration of each request.
	EnableDurMetric bool

	// EnableDurMetric whether to enable a metric that tells the number of current in-flight requests.
	EnableInFlightMetric bool
}

// HTTPCfg has the default configuration options for HTTPRecorder.
var HTTPCfg = HTTPRecorderConfig{
	EnableTotalMetric:    true,
	EnableDurMetric:      true,
	EnableInFlightMetric: true,
}

func (c HTTPRecorderConfig) SetNamespace(v string) HTTPRecorderConfig {
	c.Namespace = v

	return c
}

// HTTPLabels will contain HTTP label values for each added metric. Not all labels apply to all metrics, read the
// documentation in each metric method to find out which labels are available for that metric.
type HTTPLabels struct {
	// Method should be the HTTP method in the HTTP request.
	Method string
	// Code should be the HTTP status code in the HTTP response. If there is no response, the Code should be 0.
	Code int
	// Path is the request URL's path. Should not contain the query string and ideally it should only be the route
	// definition. For example `/users/{ID}` instead of `/users/100`.
	Path string
}

// HTTPRecorder is a recorder of HTTP metrics for prometheus. Use NewHTTPRecorder to initialize it.
type HTTPRecorder struct {
	cfg         HTTPRecorderConfig
	mp          metric.MeterProvider
	reqTotal    metric.Int64Counter
	reqDuration metric.Float64Histogram
	reqInFlight metric.Int64UpDownCounter
}

// NewHTTPRecorder creates a new HTTPRecorder. Calling this function will automatically register the new metrics to reg.
// If meter provider is nil, it will use otel's global provider. More information about configuration options in cfg can be
// found in documentation for HTTPRecorderConfig.
func NewHTTPRecorder(cfg HTTPRecorderConfig, meterProvider metric.MeterProvider) *HTTPRecorder {
	if meterProvider == nil {
		meterProvider = otel.GetMeterProvider()
	}

	m := HTTPRecorder{
		cfg: cfg,
		mp:  meterProvider,
	}

	if err := m.register(); err != nil {
		// possible errors here include duplicate metric or same metrics with inconsistent labels or help strings. It is
		// unlikely that it will happen if not by mistake. None the less we would like to know if such case occurs, hence
		// a panic
		panic(err)
	}

	return &m
}

func (h *HTTPRecorder) namespacedValue(v string) string {
	if h.cfg.Namespace != "" {
		return h.cfg.Namespace + "_" + v
	}

	return v
}

func (h *HTTPRecorder) register() error {
	meter := h.mp.Meter("")

	if h.cfg.EnableTotalMetric {
		reqTotal, err := meter.Int64Counter(
			h.namespacedValue("http_requests_total"),
			metric.WithDescription("The total number of requests"),
		)
		if err != nil {
			return fmt.Errorf("meter %s cannot set; %w", "http_requests_total", err)
		}

		h.reqTotal = reqTotal
	}

	if h.cfg.EnableDurMetric {
		reqDuration, err := meter.Float64Histogram(
			h.namespacedValue("request_duration_seconds"),
			metric.WithDescription("The total duration of a request in seconds"),
		)
		if err != nil {
			return fmt.Errorf("meter %s cannot set; %w", "request_duration_seconds", err)
		}

		h.reqDuration = reqDuration
	}

	if h.cfg.EnableInFlightMetric {
		reqInFlight, err := meter.Int64UpDownCounter(
			h.namespacedValue("requests_inflight_total"),
			metric.WithDescription("The current number of in-flight requests"),
		)
		if err != nil {
			return fmt.Errorf("meter %s cannot set; %w", "requests_inflight_total", err)
		}

		h.reqInFlight = reqInFlight
	}

	return nil
}

// AddRequestToTotal adds 1 to the total number of requests. All labels should be specified.
func (h *HTTPRecorder) AddRequestToTotal(ctx context.Context, values HTTPLabels) {
	if h.reqTotal == nil {
		return
	}

	h.reqTotal.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("method", values.Method),
			attribute.Int("code", values.Code),
		),
	)
}

// AddRequestDuration registers a request along with its duration. All labels should be specified.
func (h *HTTPRecorder) AddRequestDuration(ctx context.Context, duration time.Duration, values HTTPLabels) {
	if h.reqDuration == nil {
		return
	}

	h.reqDuration.Record(
		ctx, duration.Seconds(),
		metric.WithAttributes(
			attribute.String("method", values.Method),
			attribute.String("path", values.Path),
			attribute.Int("code", values.Code),
		),
	)
}

// AddInFlightRequest Adds 1 to the number of current in-flight requests. All labels should be specified except for
// `Code`, as it will just be ignored. To remove a request use RemInFlightRequest.
func (h *HTTPRecorder) AddInFlightRequest(ctx context.Context, values HTTPLabels) {
	if h.reqInFlight == nil {
		return
	}

	h.reqInFlight.Add(
		ctx,
		1,
		metric.WithAttributes(attribute.String("method", values.Method), attribute.String("path", values.Path)),
	)
}

// RemInFlightRequest Remove 1 from the number of current in-flight requests. All labels should be specified except
// for `Code`, as it will just be ignored. Labels should match the ones passed to the equivalent AddInFlightRequest call.
func (h *HTTPRecorder) RemInFlightRequest(ctx context.Context, values HTTPLabels) {
	if h.reqInFlight == nil {
		return
	}

	h.reqInFlight.Add(
		ctx,
		-1,
		metric.WithAttributes(attribute.String("method", values.Method), attribute.String("path", values.Path)),
	)
}
