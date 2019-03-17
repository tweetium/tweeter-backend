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
	AfterEach(func() {
		db.RollbackTransactionForTests()
	})

	BeforeEach(func() {
		// Transaction should be before all other actions
		db.BeginTransactionForTests()
	})

	Describe("created via Create", func() {
		Context("with a valid email and password", func() {
			var (
				user User
			)

			BeforeEach(func() {
				var createErr error
				user, createErr = Create("darren@onesignal.com", "password")
				Expect(createErr).NotTo(HaveOccurred())
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

			It("should err subsequent Creates with same email", func() {
				_, err := Create("darren@onesignal.com", "someotherpassword")
				Expect(err).To(Equal(ErrUserEmailAlreadyExists))
			})
		})

		Context("with password too short", func() {
			It("should error with ErrPasswordTooShort", func() {
				_, createErr := Create("darren@onesignal.com", "pass")
				Expect(createErr).To(Equal(ErrPasswordTooShort))
			})
		})
	})

	Describe("getting via Get", func() {
		var (
			createUser User
		)

		JustBeforeEach(func() {
			// create pre-existing user
			var err error
			createUser, err = Create("darren.tsung@gmail", "password")
			Expect(err).ToNot(HaveOccurred())
		})

		Context("getting createUser.ID", func() {
			It("returns same user as createUser", func() {
				Expect(Get(createUser.ID)).To(Equal(createUser))
			})
		})

		Context("getting non-existent ID", func() {
			It("errors with ErrUserNotFound", func() {
				_, err := Get(999)
				Expect(err).To(Equal(ErrUserNotFound))
			})
		})
	})

	Describe("parsing IDs via ParseID", func() {
		Context("positive integer string", func() {
			It("works and returns ID", func() {
				Expect(ParseID("123")).To(Equal(ID(123)))
			})
		})

		Context("negative integer string", func() {
			It("errors with ErrUserIDNotValid", func() {
				_, err := ParseID("-123")
				Expect(err).To(Equal(ErrUserIDNotValid))
			})
		})

		Context("invalid integer string", func() {
			It("errors with ErrUserIDNotValid", func() {
				_, err := ParseID("abc-123")
				Expect(err).To(Equal(ErrUserIDNotValid))
			})
		})
	})
})
