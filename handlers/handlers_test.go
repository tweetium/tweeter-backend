package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"tweeter/handlers"
	. "tweeter/handlers/testutil"
	. "tweeter/testutil"
)

func TestGetRouter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GetRouter Suite")
}

var _ = Describe("GetRouter", func() {
	var (
		server   *httptest.Server
		request  RequestArgs
		response *http.Response
	)

	AfterEach(func() {
		server.Close()
	})

	JustBeforeEach(func() {
		server = httptest.NewServer(handlers.GetRouter())
		response = MustSendRequest(server, request)
	})

	Context("When sending to healthcheck endpoint", func() {
		BeforeEach(func() {
			request = RequestArgs{
				Method:   http.MethodGet,
				Endpoint: "/healthcheck",
				RawBody:  StrPtr(""),
			}
		})

		It("has a success response", func() {
			MustReadSuccessData(response, nil)
		})
	})
})
