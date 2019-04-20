package jwt_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	usersJWT "tweeter/handlers/endpoints/users/jwt"
	"tweeter/jwtsecrets"
)

func TestUsersJWT(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Users JWT Suite")
}

var _ = Describe("Users JWT", func() {
	var (
		signedToken    string
		expirationTime time.Time
		generateError  error

		parsedClaims usersJWT.Claims
		parsedError  error
	)

	var originalClaims = usersJWT.Claims{
		UserID: 10,
	}

	BeforeEach(func() {
		jwtsecrets.Init(
			map[string]string{"1": "03ad766e-1ef5-4019-98e5-e65beb286ae3"},
			"1", // the current key
		)

		signedToken, generateError = usersJWT.GenerateTokenWithExpiration(originalClaims, expirationTime)
		parsedClaims, parsedError = usersJWT.ParseToken(signedToken)
	})

	Context("signed token in future", func() {
		JustBeforeEach(func() {
			expirationTime = time.Now().Add(time.Minute * 30)
		})

		It("doesn't error on generating", func() {
			Expect(generateError).NotTo(HaveOccurred())
		})

		It("doesn't error on parsing", func() {
			Expect(parsedError).NotTo(HaveOccurred())
		})

		It("parses claims correctly", func() {
			Expect(parsedClaims).To(Equal(originalClaims))
		})

		Context("after secret used is removed from list", func() {
			JustBeforeEach(func() {
				jwtsecrets.Init(
					map[string]string{"2": "bde63f2d-be16-4493-810a-2ff897f1647b"},
					"2", // the current key
				)

				// reparse the token with new secrets
				parsedClaims, parsedError = usersJWT.ParseToken(signedToken)
			})

			It("errors on parsing", func() {
				Expect(parsedError).To(HaveOccurred())
			})
		})

		Context("when secret is retired and new secret is used", func() {
			JustBeforeEach(func() {
				jwtsecrets.Init(
					map[string]string{
						"1": "03ad766e-1ef5-4019-98e5-e65beb286ae3",
						"2": "bde63f2d-be16-4493-810a-2ff897f1647b",
					},
					"2", // the current key
				)

				// reparse the token with new secrets
				parsedClaims, parsedError = usersJWT.ParseToken(signedToken)
			})

			It("doesn't error on parsing", func() {
				Expect(parsedError).NotTo(HaveOccurred())
			})

			It("parses claims correctly", func() {
				Expect(parsedClaims).To(Equal(originalClaims))
			})
		})
	})

	Context("signed token in past", func() {
		JustBeforeEach(func() {
			expirationTime = time.Now().Add(-time.Minute * 30)
		})

		It("doesn't error on generating", func() {
			Expect(generateError).NotTo(HaveOccurred())
		})

		It("errors on parsing with ErrTokenExpired", func() {
			Expect(parsedError).To(Equal(jwtsecrets.ErrTokenExpired))
		})
	})
})
