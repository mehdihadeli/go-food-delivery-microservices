package otelmetrics

import (
	"go.opentelemetry.io/otel/sdk/metric"
)

// https://opentelemetry.io/docs/instrumentation/go/manual/#registering-views
func GetViews() []metric.View {
	customBucketView := metric.NewView(
		metric.Instrument{
			Name: "*request_duration_seconds",
		},
		metric.Stream{
			Aggregation: metric.AggregationExplicitBucketHistogram{
				Boundaries: []float64{
					.005,
					.01,
					.025,
					.05,
					.1,
					.25,
					.5,
					1,
					2.5,
					5,
					10,
				},
			},
		},
	)

	return []metric.View{customBucketView}
}
