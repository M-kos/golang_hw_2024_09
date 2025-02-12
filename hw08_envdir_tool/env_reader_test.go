package main

import (
	"testing"

	"github.com/stretchr/testify/require" //nolint:all
)

func TestReadDir(t *testing.T) {
	cmd := []string{"ls", "-la"}

	t.Run("success test", func(t *testing.T) {
		code := RunCmd(cmd, Environment{})

		require.Equal(t, 0, code)
	})

	t.Run("success case with env", func(t *testing.T) {
		code := RunCmd(cmd, Environment{
			"FOO": EnvValue{
				Value: "BAR",
			},
		})

		require.Equal(t, 0, code)
	})

	t.Run("Exec err case", func(t *testing.T) {
		cmd[1] = "kk"

		code := RunCmd(cmd, Environment{})

		require.Equal(t, 2, code)
	})
}
