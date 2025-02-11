package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require" //nolint:all
)

const (
	testDir = "./testdata/env"
)

func TestRunCmd(t *testing.T) {
	t.Run("successful test", func(t *testing.T) {
		expectedEnv := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: " ", NeedRemove: false},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}

		env, err := ReadDir(testDir)

		fmt.Println("ENV: ", env)

		require.NoError(t, err)
		require.Equal(t, expectedEnv, env)
	})
	t.Run("test with empty dir", func(t *testing.T) {
		tempDir := t.TempDir()

		res, err := ReadDir(tempDir)
		require.NoError(t, err)
		require.Len(t, res, 0)
	})
}
