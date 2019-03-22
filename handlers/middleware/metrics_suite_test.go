package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"

	"tweeter/handlers/middleware"
	"tweeter/metrics"
	. "tweeter/testutil"
)

func TestMetricsMiddleware(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Metrics Middleware Endpoint Suite")
}

var _ = Describe("Metrics Middleware", func() {
	var (
		endpointName = "TestEndpoint"
		server       *httptest.Server
	)

	JustBeforeEach(func() {
		handler := middleware.Metrics(endpointName, func(w http.ResponseWriter, r *http.Request) {})
		server = httptest.NewServer(handler)
	})

	AfterEach(func() {
		server.Close()
	})

	It("records to APIReceived after receiving request", func() {
		collector := metrics.APIReceived.
			With(prometheus.Labels{
				"endpointName": "TestEndpoint",
			})

		diff := PrometheusCollectorDiff(
			collector,
			func() { _, _ = server.Client().Get(server.URL) },
		)

		Expect(diff).To(Equal(1.0))
	})

})
