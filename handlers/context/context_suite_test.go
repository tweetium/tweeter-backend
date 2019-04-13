package context_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/prometheus/client_golang/prometheus"

	"tweeter/handlers/context"
	"tweeter/handlers/responses"
	"tweeter/metrics"
	. "tweeter/testutil"
)

func TestContext(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Context Endpoint Suite")
}

var _ = Describe("Context", func() {
	var (
		endpointName = "TestEndpoint"
		context      = context.New(
			endpointName,
			httptest.NewRecorder(),
			MustNewRequest(
				http.MethodGet,
				"/test-endpoint/",
				strings.NewReader("Hello world!"),
			),
		)
		statusCode = 200
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
				func() { context.RenderResponse(statusCode, response) },
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
				responses.Error{
					Title:  "Second Error Title",
					Detail: "Second Detail Blah Blah Blah",
				},
			}
		)

		JustBeforeEach(func() {
			context.RenderErrorResponse(statusCode, errors...)
		})

		It("records metrics to APIResponses with correct labels", func() {
			collector := metrics.APIResponses.
				With(prometheus.Labels{
					"endpointName": "TestEndpoint",
					"statusCode":   "200",
				})

			diff := PrometheusCollectorDiff(
				collector,
				func() { context.RenderErrorResponse(statusCode, errors...) },
			)

			Expect(diff).To(Equal(1.0))
		})

		It("records the error title to APIResponseErrors", func() {
			collector := metrics.APIResponseErrors.
				With(prometheus.Labels{
					"endpointName": "TestEndpoint",
					"errorTitle":   "Test Error Title",
				})

			diff := PrometheusCollectorDiff(
				collector,
				func() { context.RenderErrorResponse(statusCode, errors...) },
			)

			Expect(diff).To(Equal(1.0))
		})

		It("records the second error title to APIResponseErrors", func() {
			collector := metrics.APIResponseErrors.
				With(prometheus.Labels{
					"endpointName": "TestEndpoint",
					"errorTitle":   "Second Error Title",
				})

			diff := PrometheusCollectorDiff(
				collector,
				func() { context.RenderErrorResponse(statusCode, errors...) },
			)

			Expect(diff).To(Equal(1.0))
		})
	})
})
