package metrics_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config"
	customEcho "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/customecho"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/external/fxlog"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/zap"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/metrics"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHealth(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "/health suite")
}

var _ = Describe("/", Ordered, func() {
	var (
		url string
		err error
		res *http.Response
	)

	BeforeAll(func() {
		var cfg *metrics.MetricsOptions

		fxtest.New(
			GinkgoT(),
			zap.Module,
			fxlog.FxLogger,
			config.Module,
			customEcho.Module,

			metrics.Module,

			fx.Populate(&cfg),
		).RequireStart()

		url = fmt.Sprintf("http://%s:%s/metrics", cfg.Host, cfg.Port)
	})

	BeforeEach(func() {
		res, err = http.Get(url)
	})
	It("returns status OK", func() {
		Expect(err).To(BeNil())
		Expect(res.StatusCode).To(Equal(http.StatusOK))
	})

	It("returns how many requests were made", func() {
		b, err := io.ReadAll(res.Body)
		Expect(err).To(BeNil())

		Expect(
			b,
		).To(ContainSubstring(`promhttp_metric_handler_requests_total{code="200"} 1`))
	})
})
