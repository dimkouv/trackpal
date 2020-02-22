package envlib

import (
	"os"

	"github.com/sirupsen/logrus"
)

const EnvPostgresDSN = "POSTGRES_DSN"

func GetEnvOrPanic(envVarKey string) string {
	v := os.Getenv(envVarKey)
	if v == "" {
		panic("environment variable '" + envVarKey + "' not found")
	}
	return v
}

func GetEnvOrDefault(envVarKey, defaultVal string) string {
	v := os.Getenv(envVarKey)
	if v == "" {
		logrus.Warnf("Environment variable with key='%s' not found. Fallback value='%s'", envVarKey, defaultVal)
		return defaultVal
	}
	return v
}
