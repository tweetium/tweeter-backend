package login_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"tweeter/db"
	"tweeter/db/models/user"
	"tweeter/handlers/endpoints/users"
	usersJWT "tweeter/handlers/endpoints/users/jwt"
	usersLogin "tweeter/handlers/endpoints/users/login"
	"tweeter/handlers/responses"
	. "tweeter/handlers/testutil"
)

func TestLoginEndpoint(t *testing.T) {
	db.InitForTests()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Users#login Endpoint Suite")
}

var _ = Describe("Users#login Endpoint", func() {
	var (
		server   *httptest.Server
		request  RequestArgs
		response *http.Response

		existingUser user.User
	)

	var existingUserEmail = "darren.t@gmail.com"
	var existingUserPassword = "test_password"

	var sendRequest = func(request RequestArgs) *http.Response {
		return MustSendRequest(server, request)
	}

	AfterEach(func() {
		server.Close()
		db.RollbackTransactionForTests()
	})

	BeforeEach(func() {
		// Transaction should be before all other actions
		db.BeginTransactionForTests()

		var err error
		existingUser, err = user.Create(existingUserEmail, existingUserPassword)
		Expect(err).NotTo(HaveOccurred())
	})

	JustBeforeEach(func() {
		r := mux.NewRouter()
		usersLogin.Endpoint.Attach(r)
		server = httptest.NewServer(r)
		response = sendRequest(request)
	})

	var requestArgs = func(email, password string) RequestArgs {
		return RequestArgs{
			Method:   http.MethodPost,
			Endpoint: "/api/users/login",
			JSONBody: map[string]interface{}{
				"email":    email,
				"password": password,
			},
		}
	}

	var successfulRequestArgs = func() RequestArgs {
		return requestArgs(existingUserEmail, existingUserPassword)
	}

	Context("without secrets initialized", func() {
		BeforeEach(func() { usersJWT.ClearSecretsMap() })

		Context("with a valid email and password", func() {
			BeforeEach(func() {
				request = successfulRequestArgs()
			})

			It("errors with internal error", func() {
				errors := MustReadErrors(response)
				Expect(errors).To(Equal([]responses.Error{responses.ErrInternalError}))
			})
		})
	})

	Context("with secrets initialized", func() {
		BeforeEach(func() {
			usersJWT.InitializeWithSecretsMap(
				map[string]string{"1": "03ad766e-1ef5-4019-98e5-e65beb286ae3"},
				"1", // the current key
			)
		})

		Context("with a valid email and password", func() {
			BeforeEach(func() {
				request = successfulRequestArgs()
			})

			It("has a success response", func() {
				MustReadSuccessData(response, nil)
			})

			It("sets the usersJWT cookie correctly", func() {
				cookies := response.Cookies()
				Expect(len(cookies)).To(Equal(1))

				expirationTime := time.Now().AddDate(1 /* year */, 0, 0)
				// Add some buffer for the test running slow
				expirationTime = expirationTime.Add(10 * time.Second)

				cookie := cookies[0]
				Expect(cookie.Name).To(Equal(usersJWT.CookieName))
				Expect(cookie.Expires).Should(BeTemporally("<", expirationTime))

				claims, err := usersJWT.ParseToken(cookie.Value)
				Expect(err).NotTo(HaveOccurred())

				Expect(claims.UserID).To(Equal(existingUser.ID))
			})
		})

		Context("with an invalid email", func() {
			BeforeEach(func() {
				request = requestArgs("invalid@gmail.com", existingUserPassword)
			})

			It("errors with invalid credentials error", func() {
				errors := MustReadErrors(response)
				Expect(errors).To(Equal([]responses.Error{users.ErrInvalidCredentials}))
			})
		})

		Context("with an invalid password", func() {
			BeforeEach(func() {
				request = requestArgs(existingUserEmail, "invalid_password")
			})

			It("errors with invalid credentials error", func() {
				errors := MustReadErrors(response)
				Expect(errors).To(Equal([]responses.Error{users.ErrInvalidCredentials}))
			})
		})
	})
})
