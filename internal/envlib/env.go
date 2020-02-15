package envlib

import "os"

const EnvPostgresDSN = "POSTGRES_DSN"

func GetEnvOrPanic(envVarKey string) string {
	v := os.Getenv(envVarKey)
	if v == "" {
		panic("environment variable '" + envVarKey + "' not found")
	}
	return v
}
