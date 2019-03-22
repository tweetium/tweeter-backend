package render_test

import (
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/prometheus/client_golang/prometheus"

	"tweeter/handlers/render"
	"tweeter/handlers/responses"
	"tweeter/metrics"
	. "tweeter/testutil"
)

func TestRender(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Render Endpoint Suite")
}

var _ = Describe("Render", func() {
	var (
		endpointName   = "TestEndpoint"
		responseWriter = httptest.NewRecorder()
		statusCode     = 200
	)

	Context("with Response call", func() {
		var (
			response = responses.NewSuccessResponse(nil)
		)

		It("records metrics to APIResponses with correct labels", func() {
			collector := metrics.APIResponses.
				With(prometheus.Labels{
					"endpointName": "TestEndpoint",
					"statusCode":   "200",
				})

			diff := PrometheusCollectorDiff(
				collector,
				func() { render.Response(endpointName, responseWriter, statusCode, response) },
			)

			Expect(diff).To(Equal(1.0))
		})
	})

	Context("with ErrorResponse call", func() {
		var (
			errors = []responses.Error{
				responses.Error{
					Title:  "Test Error Title",
					Detail: "Test Detail Blah Blah Blah",
				},
			}
		)

		JustBeforeEach(func() {
			render.ErrorResponse(endpointName, responseWriter, statusCode, errors...)
		})

		It("records metrics to APIResponses with correct labels", func() {
			collector := metrics.APIResponses.
				With(prometheus.Labels{
					"endpointName": "TestEndpoint",
					"statusCode":   "200",
				})

			diff := PrometheusCollectorDiff(
				collector,
				func() { render.ErrorResponse(endpointName, responseWriter, statusCode, errors...) },
			)

			Expect(diff).To(Equal(1.0))
		})

		It("records metrics to APIResponseErrors with correct labels", func() {
			collector := metrics.APIResponseErrors.
				With(prometheus.Labels{
					"endpointName": "TestEndpoint",
					"errorTitle":   "Test Error Title",
				})

			diff := PrometheusCollectorDiff(
				collector,
				func() { render.ErrorResponse(endpointName, responseWriter, statusCode, errors...) },
			)

			Expect(diff).To(Equal(1.0))
		})
	})
})
