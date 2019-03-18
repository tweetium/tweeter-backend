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
		request     *http.Request
		response    *http.Response
		responseErr error
	)

	AfterEach(func() {
		db.RollbackTransactionForTests()
	})

	BeforeEach(func() {
		// Transaction should be before all other actions
		db.BeginTransactionForTests()
	})

	JustBeforeEach(func() {
		server := httptest.NewServer(nil)
		request.URL = MustURLParse(server.URL + request.URL.Path)
		response, responseErr = server.Client().Do(request)
		server.Close()
	})

	Describe("creating users via POST", func() {
		var basicSuccessfulRequest = func() *http.Request {
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
				request = basicSuccessfulRequest()
			})

			It("should not have errored", func() {
				Expect(responseErr).NotTo(HaveOccurred())
			})

			It("has a success response with non-zero ID", func() {
				idData := struct{ ID user.ID }{}
				MustReadSuccessData(response, idData)
				Expect(idData.ID).NotTo(Equal(0))
			})
		})
	})
})
