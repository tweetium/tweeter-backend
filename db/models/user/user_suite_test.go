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

			It("has created and modified time within time limit", func() {
				// Maximum difference allowed due to latency with tests / db
				var MaxTimeDiff = 1 * time.Second

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

	Describe("finding users via Find", func() {
		var (
			createUser User

			findBy FindBy

			findUser  User
			findError error
		)

		BeforeEach(func() {
			// create pre-existing user
			var err error
			createUser, err = Create("darren.tsung@gmail.com", "password")
			Expect(err).ToNot(HaveOccurred())
		})

		JustBeforeEach(func() {
			findUser, findError = Find(findBy)
		})

		Context("using FindByID", func() {
			Context("with createUser.ID", func() {
				BeforeEach(func() { findBy = FindByID{ID: createUser.ID} })

				It("returns same user as createUser", func() {
					Expect(findUser).To(Equal(createUser))
				})

				It("returns no error", func() {
					Expect(findError).ToNot(HaveOccurred())
				})
			})

			Context("with non-existent ID", func() {
				BeforeEach(func() { findBy = FindByID{ID: 999} })

				It("errors with ErrUserNotFound", func() {
					Expect(findError).To(Equal(ErrUserNotFound))
				})
			})
		})

		Context("using FindByEmail", func() {
			Context("with createUser.Email", func() {
				BeforeEach(func() { findBy = FindByEmail{Email: createUser.Email} })

				It("returns same user as createUser", func() {
					Expect(findUser).To(Equal(createUser))
				})

				It("returns no error", func() {
					Expect(findError).ToNot(HaveOccurred())
				})
			})

			Context("with non-existent ID", func() {
				BeforeEach(func() { findBy = FindByEmail{Email: "invalid@gmail.com"} })

				It("errors with ErrUserNotFound", func() {
					Expect(findError).To(Equal(ErrUserNotFound))
				})
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

	Describe("validating users via ComparePassword", func() {
		var (
			createUser            User
			createUserRawPassword string

			comparePassword string

			passwordMatches bool
		)

		BeforeEach(func() {
			// create pre-existing user
			var err error
			// need to use raw password for validate because createUser.Password is hashed :)
			createUserRawPassword = "password"
			createUser, err = Create("darren.tsung@gmail.com", createUserRawPassword)
			Expect(err).ToNot(HaveOccurred())
		})

		JustBeforeEach(func() {
			passwordMatches = createUser.ComparePassword(comparePassword)
		})

		Context("with correct password", func() {
			BeforeEach(func() {
				comparePassword = createUserRawPassword
			})

			It("returns true (password matches)", func() {
				Expect(passwordMatches).To(Equal(true))
			})
		})

		Context("with incorrect password", func() {
			BeforeEach(func() {
				comparePassword = "wrongpassword"
			})

			It("returns false (password doesn't match)", func() {
				Expect(passwordMatches).To(Equal(false))
			})
		})
	})
})
