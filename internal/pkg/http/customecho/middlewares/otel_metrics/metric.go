package otelmetrics

// ref:https://github.com/open-telemetry/opentelemetry-go/blob/main/example/prometheus/main.go
// https://github.com/labstack/echo-contrib/blob/master/prometheus/prometheus.go
// https://github.com/worldline-go/tell/tree/main/metric/metricecho
// https://opentelemetry.io/docs/instrumentation/go/manual/#metrics
// https://opentelemetry.io/docs/specs/otel/metrics/semantic_conventions/http-metrics/

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// HTTPLabels will contain HTTP label values for each added metric. Not all labels apply to all metrics, read the
// documentation in each metric method to find out which labels are available for that metric.
type HTTPLabels struct {
	// Method should be the HTTP method in the HTTP request.
	Method string
	// Code should be the HTTP status code in the HTTP response. If there is no response, the Code should be 0.
	Code int
	// Path is the request URL's path. Should not contain the query string, and ideally it should only be the route
	// definition. For example `/users/{ID}` instead of `/users/100`.
	Path string
	Host string
}

// HTTPMetricsRecorder is a recorder of HTTP metrics for prometheus. Use NewHTTPMetricsRecorder to initialize it.
type HTTPMetricsRecorder struct {
	cfg            config
	mp             metric.MeterProvider
	meter          metric.Meter
	reqTotal       metric.Int64Counter
	reqDuration    metric.Float64Histogram
	reqInFlight    metric.Int64UpDownCounter
	resSize        metric.Int64Histogram
	reqSize        metric.Int64Histogram
	errorCounter   metric.Float64Counter
	successCounter metric.Float64Counter
}

// NewHTTPMetricsRecorder creates a new HTTPMetricsRecorder. Calling this function will automatically register the new metrics to reg.
func NewHTTPMetricsRecorder(cfg config) *HTTPMetricsRecorder {
	// Meter can be a global/package variable.
	meter := cfg.metricsProvider.Meter(cfg.instrumentationName)

	m := HTTPMetricsRecorder{
		cfg:   cfg,
		mp:    cfg.metricsProvider,
		meter: meter,
	}

	if err := m.register(); err != nil {
		// possible errors here include duplicate metric or same metrics with inconsistent labels or help strings. It is
		// unlikely that it will happen if not by mistake. Nonetheless, we would like to know if such case occurs, hence
		// a panic
		panic(err)
	}

	return &m
}

func (h *HTTPMetricsRecorder) namespacedValue(v string) string {
	if h.cfg.namespace != "" {
		return h.cfg.namespace + "_" + v
	}

	return v
}

func (h *HTTPMetricsRecorder) register() error {
	// https://opentelemetry.io/docs/specs/otel/metrics/semantic_conventions/http-metrics/#http-server
	errorCounter, err := h.meter.Float64Counter(
		"http.server.total_error_request",
		metric.WithUnit("count"),
		metric.WithDescription("The total number of error http requests"),
	)
	if err != nil {
		return fmt.Errorf(
			"meter %s cannot set; %w",
			"http.server.total_error_request",
			err,
		)
	}

	h.errorCounter = errorCounter

	successCounter, err := h.meter.Float64Counter(
		"http.server.total_success_request",
		metric.WithUnit("count"),
		metric.WithDescription("The total number of success http requests"),
	)
	if err != nil {
		return fmt.Errorf(
			"meter %s cannot set; %w",
			"http.server.total_success_request",
			err,
		)
	}

	h.successCounter = successCounter

	if h.cfg.enableTotalMetric {
		reqTotal, err := h.meter.Int64Counter(
			h.namespacedValue("http.server.total_request"),
			metric.WithUnit("count"),
			metric.WithDescription("The total number of requests"),
		)
		if err != nil {
			return fmt.Errorf(
				"meter %s cannot set; %w",
				"http.server.total_request",
				err,
			)
		}

		h.reqTotal = reqTotal
	}

	if h.cfg.enableDurMetric {
		reqDuration, err := h.meter.Float64Histogram(
			h.namespacedValue("http.server.duration"),
			metric.WithUnit("s"), // Specify the unit as "seconds"
			metric.WithDescription(
				"The total duration of a request in seconds",
			),
		)
		if err != nil {
			return fmt.Errorf(
				"meter %s cannot set; %w",
				"http.server.duration",
				err,
			)
		}

		h.reqDuration = reqDuration
	}

	if h.cfg.enableInFlightMetric {
		reqInFlight, err := h.meter.Int64UpDownCounter(
			h.namespacedValue("http.server.request_inflight_total"),
			metric.WithUnit("count"),
			metric.WithDescription("The current number of in-flight requests"),
		)
		if err != nil {
			return fmt.Errorf(
				"meter %s cannot set; %w",
				"http.server.request_inflight_total",
				err,
			)
		}

		h.reqInFlight = reqInFlight
	}

	resSize, err := h.meter.Int64Histogram(
		h.namespacedValue("http.server.response.size"),
		metric.WithUnit("bytes"),
		metric.WithDescription("The HTTP response sizes in bytes."),
	)
	if err != nil {
		return fmt.Errorf(
			"meter %s cannot set; %w",
			"http.server.response.size",
			err,
		)
	}

	h.resSize = resSize

	reqSize, err := h.meter.Int64Histogram(
		h.namespacedValue("http.server.request.size"),
		metric.WithUnit("bytes"),
		metric.WithDescription("The HTTP request sizes in bytes."),
	)
	if err != nil {
		return fmt.Errorf(
			"meter %s cannot set; %w",
			"http.server.request.size",
			err,
		)
	}

	h.reqSize = reqSize

	return nil
}

// AddRequestToTotal adds 1 to the total number of requests. All labels should be specified.
func (h *HTTPMetricsRecorder) AddRequestToTotal(
	ctx context.Context,
	values HTTPLabels,
) {
	if h.reqTotal == nil {
		return
	}

	h.reqTotal.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("method", values.Method),
			attribute.Int("code", values.Code),
			attribute.String("type", "Http"),
		),
	)
}

// AddRequestDuration registers a request along with its duration. All labels should be specified.
func (h *HTTPMetricsRecorder) AddRequestDuration(
	ctx context.Context,
	duration time.Duration,
	values HTTPLabels,
) {
	if h.reqDuration == nil {
		return
	}

	h.reqDuration.Record(
		ctx, duration.Seconds(),
		metric.WithAttributes(
			attribute.String("method", values.Method),
			attribute.String("host", values.Host),
			attribute.String("path", values.Path),
			attribute.Int("code", values.Code),
			attribute.String("type", "Http"),
		),
	)
}

// AddInFlightRequest Adds 1 to the number of current in-flight requests. All labels should be specified except for
// `Code`, as it will just be ignored. To remove a request use RemInFlightRequest.
func (h *HTTPMetricsRecorder) AddInFlightRequest(
	ctx context.Context,
	values HTTPLabels,
) {
	if h.reqInFlight == nil {
		return
	}

	h.reqInFlight.Add(
		ctx,
		1,
		metric.WithAttributes(
			attribute.String("method", values.Method),
			attribute.String("path", values.Path),
			attribute.String("type", "Http"),
		),
	)
}

func (h *HTTPMetricsRecorder) AddRequestError(
	ctx context.Context,
	values HTTPLabels,
) {
	if h.errorCounter == nil {
		return
	}

	h.errorCounter.Add(
		ctx,
		1,
		metric.WithAttributes(
			attribute.String("method", values.Method),
			attribute.String("path", values.Path),
			attribute.String("type", "Http"),
			attribute.Int("code", values.Code),
		),
	)
}

func (h *HTTPMetricsRecorder) AddRequestSuccess(
	ctx context.Context,
	values HTTPLabels,
) {
	if h.successCounter == nil {
		return
	}

	h.successCounter.Add(
		ctx,
		1,
		metric.WithAttributes(
			attribute.String("method", values.Method),
			attribute.String("path", values.Path),
			attribute.String("type", "Http"),
			attribute.Int("code", values.Code),
		),
	)
}

// RemInFlightRequest Remove 1 from the number of current in-flight requests. All labels should be specified except
// for `Code`, as it will just be ignored. Labels should match the ones passed to the equivalent AddInFlightRequest call.
func (h *HTTPMetricsRecorder) RemInFlightRequest(
	ctx context.Context,
	values HTTPLabels,
) {
	if h.reqInFlight == nil {
		return
	}

	h.reqInFlight.Add(
		ctx,
		-1,
		metric.WithAttributes(
			attribute.String("method", values.Method),
			attribute.String("path", values.Path),
			attribute.String("type", "Http"),
		),
	)
}

func (h *HTTPMetricsRecorder) AddRequestSize(
	ctx context.Context,
	request *http.Request,
	values HTTPLabels,
) {
	if h.reqSize == nil {
		return
	}

	size := computeApproximateRequestSize(request)
	h.reqSize.Record(ctx, int64(size), metric.WithAttributes(
		attribute.String("method", values.Method),
		attribute.String("path", values.Path),
		attribute.String("type", "Http"),
		attribute.String("host", values.Host),
		attribute.Int("code", values.Code),
	))
}

func (h *HTTPMetricsRecorder) AddResponseSize(
	ctx context.Context,
	response *echo.Response,
	values HTTPLabels,
) {
	if h.resSize == nil {
		return
	}

	size := response.Size
	h.resSize.Record(ctx, size, metric.WithAttributes(
		attribute.String("method", values.Method),
		attribute.String("path", values.Path),
		attribute.String("type", "Http"),
		attribute.String("host", values.Host),
		attribute.Int("code", values.Code),
	))
}

func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.Path)
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}

	return s
}
