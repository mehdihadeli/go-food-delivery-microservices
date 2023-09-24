package metricecho

import (
	"github.com/worldline-go/tell/tglobal"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
)

func GetViews() []metric.View {
	customBucketView := metric.NewView(
		metric.Instrument{
			Name: "*request_duration_seconds",
		},
		metric.Stream{
			Aggregation: aggregation.ExplicitBucketHistogram{
				Boundaries: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
		},
	)

	return []metric.View{customBucketView}
}

func init() {
	tglobal.MetricViews.Add("echo", GetViews())
}
