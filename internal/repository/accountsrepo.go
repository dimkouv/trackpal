package repository

import "github.com/dimkouv/trackpal/internal/models"

// UserAccountRepository contains methods for operations on user accounts
type UserAccountRepository interface {
	// SaveNewUser can be used for storing a new user account. You need to provide a password which can be later used
	// for fetching the user.
	SaveNewUser(ua models.UserAccount, password string) (*models.UserAccount, error)

	// ActivateUserAccount receives a token and an email, if the token has not expired and it matches the email
	// then the account gets activated.
	ActivateUserAccount(email, token string) error

	// GetUserByEmailAndPassword can be used for fetching a user that matches the target email and password.
	GetUserByEmailAndPassword(email, password string) (*models.UserAccount, error)

	// UpdateUser can be used for updating some information of the user. Responds
	// with whether or not a user got updated and an error
	UpdateUser(userID int64, input *UpdateUserInput) (bool, error)
}

// UpdateUserInput can be used for updating user information
type UpdateUserInput struct {
	Email           *string
	Password        *string
	FirstName       *string
	LastName        *string
	IsActive        *bool
	ActivationToken *string
}
