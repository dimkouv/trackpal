package envlib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvOrDefault(t *testing.T) {
	assert.NotEmpty(t, GetEnvOrDefault("GOPATH", "default val"))
	assert.Equal(t, "default val", GetEnvOrDefault("AN_INVALID_ENV_VAR", "default val"))
}

func TestGetEnvOrPanic(t *testing.T) {
	t.Run("should not panic of the env var exists", func(t *testing.T) {
		assert.NotEmpty(t, GetEnvOrPanic("GOPATH"))
	})

	t.Run("should raise panic if the env var does not exist", func(t *testing.T) {
		panicked := false
		defer func() {
			if err := recover(); err != nil {
				panicked = true
			}
			assert.True(t, panicked)
		}()
		GetEnvOrPanic("AN_INVALID_ENV_VAR") // this should panic
	})
}
