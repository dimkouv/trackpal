package repository

import "github.com/dimkouv/trackpal/internal/models"

// UserAccountRepository contains methods for operations on user accounts
type UserAccountRepository interface {
	// SaveNewUser can be used for storing a new user account. You need to provide a password which can be later used
	// for fetching the user.
	SaveNewUser(ua models.UserAccount, password string) (*models.UserAccount, error)

	// GetUserByEmailAndPassword can be used for fetching a user that matches the target email and password.
	GetUserByEmailAndPassword(email, password string) (*models.UserAccount, error)
}
