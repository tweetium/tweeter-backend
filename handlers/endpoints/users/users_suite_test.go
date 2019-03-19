package users_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"tweeter/db"
	"tweeter/db/models/user"
	"tweeter/handlers/endpoints/users"
	"tweeter/handlers/responses"
	. "tweeter/handlers/testutil"
	. "tweeter/testutil"
)

func TestUsersEndpoint(t *testing.T) {
	db.InitForTests()

	users.Endpoint.Attach()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Users Endpoint Suite")
}

var _ = Describe("Users Endpoint", func() {
	var (
		server   *httptest.Server
		request  *http.Request
		response *http.Response
	)

	var sendRequest = func(request *http.Request) *http.Response {
		return MustSendRequest(server, request)
	}

	AfterEach(func() {
		server.Close()
		db.RollbackTransactionForTests()
	})

	BeforeEach(func() {
		// Transaction should be before all other actions
		db.BeginTransactionForTests()
	})

	JustBeforeEach(func() {
		server = httptest.NewServer(nil)
		response = sendRequest(request)
	})

	Describe("creating users via POST", func() {
		var successfulRequest = func() *http.Request {
			req := MustNewRequest(http.MethodPost,
				"/api/v1/users",
				MustJSONMarshal(map[string]interface{}{
					"email":    "darren.tsung@gmail.com",
					"password": "password",
				}),
			)
			return req
		}

		Context("with a valid email and password", func() {
			BeforeEach(func() {
				request = successfulRequest()
			})

			It("has a success response with non-zero ID", func() {
				idData := struct{ ID user.ID }{}
				MustReadSuccessData(response, idData)
				Expect(idData.ID).NotTo(Equal(0))
			})

			It("errors for requests with the same email", func() {
				secondResponse := sendRequest(successfulRequest())
				errors := MustReadErrors(secondResponse)
				Expect(errors).To(Equal([]responses.Error{users.ErrEmailAlreadyExists("darren.tsung@gmail.com")}))
			})
		})
	})
})
