package consts

// Keys that can be used for logging
const (
	LogFieldErr  = "error"
	LogFieldBody = "body"
	LogFieldVars = "vars"
	LogFieldRepo = "repository"
)

// Keys that can be used in several places
const (
	Authorization = "Authorization"
)

type ctxKey string

// Keys that can be used as context keys
const (
	CtxUser ctxKey = "user"
)
