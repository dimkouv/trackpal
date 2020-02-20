package consts

import "github.com/dimkouv/trackpal/pkg/terror"

const (
	PQCodeForeignKeyViolation = "23503"
	PQCodeUniqueKeyViolation  = "23505"
)

var (
	ErrEnumAccountExists    = terror.New(400, "ACCOUNT_EXISTS")
	ErrEnumInvalidBody      = terror.New(400, "INVALID_BODY")
	ErrEnumInvalidEmail     = terror.New(400, "INVALID_EMAIL")
	ErrEnumInsecurePassword = terror.New(400, "INSECURE_PASSWORD")
	ErrEnumNotActivated     = terror.New(400, "NOT_ACTIVATED")
	ErrEnumNotFound         = terror.New(404, "NOT_FOUND")
	ErrInternal             = terror.New(500, "INTERNAL_ERROR")
)
