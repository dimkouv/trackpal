package conf

import (
	"time"

	"github.com/dimkouv/trackpal/internal/envlib"
)

type Argon2Conf struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// nolint: gochecknoglobals
var (
	Argon2Params = Argon2Conf{
		Memory:      32 * 1024,
		Iterations:  4,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	JWTSignBytes = []byte(envlib.GetEnvOrDefault(
		"TRACKPAL_SIGN_KEY",
		"you-should-change-this-to-something-secure",
	))

	AccessControlAllowOrigin = "*"

	// JwtTTL After this duration passes JWTs are expired
	JwtTTL = 5 * time.Minute
)

// nolint: gochecknoglobals
var (
	MailFrom = envlib.GetEnvOrDefault("TRACKPAL_MAIL_SENDER", "bot@trackpal.xyz")
)
