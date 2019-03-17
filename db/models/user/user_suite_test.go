package user_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"tweeter/db"
	. "tweeter/db/models/user"
)

func TestUser(t *testing.T) {
	db.InitForTests()

	RegisterFailHandler(Fail)
	RunSpecs(t, "User Suite")
}

var _ = Describe("User", func() {
	var (
		email     string
		password  string
		user      User
		createErr error
	)

	AfterEach(func() {
		db.RollbackTransactionForTests()
	})

	BeforeEach(func() {
		// Transaction should be before all other actions
		db.BeginTransactionForTests()
	})

	JustBeforeEach(func() {
		// Create in JustBeforeEach so email / password can be customized in BeforeEach
		user, createErr = Create(email, password)
	})

	Describe("created via Create", func() {
		Context("with a valid email and password", func() {
			BeforeEach(func() {
				email = "darren@onesignal.com"
				password = "password"
			})

			It("has same email", func() {
				Expect(user.Email).To(Equal("darren@onesignal.com"))
			})

			It("does not store plaintext password (does not prove hashed / salted)", func() {
				Expect(user.Password).NotTo(Equal("password"))
			})

			// Maximum difference allowed due to latency with tests / db
			var MaxTimeDiff = 1 * time.Second
			It("has created and modified time within time limit", func() {
				now := time.Now()
				Expect(now.Sub(user.Created)).Should(BeNumerically("<", MaxTimeDiff))
				Expect(now.Sub(user.Modified)).Should(BeNumerically("<", MaxTimeDiff))
			})

			It("should not error", func() {
				Expect(createErr).NotTo(HaveOccurred())
			})
		})
	})
})
