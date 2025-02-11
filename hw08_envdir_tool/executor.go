package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command, args := cmd[0], cmd[1:]

	exCmd := exec.Command(command, args...)
	exCmd.Stdin = os.Stdin
	exCmd.Stdout = os.Stdout
	exCmd.Stderr = os.Stderr

	for k, v := range env {
		if v.NeedRemove {
			err := os.Unsetenv(k)
			if err != nil {
				log.Println(err.Error())
				return 1
			}

			continue
		}
		exCmd.Env = append(exCmd.Env, k+"="+v.Value)
	}

	exCmd.Env = append(exCmd.Env, os.Environ()...)

	if err := exCmd.Run(); err != nil {
		log.Println(err.Error())
	}

	return exCmd.ProcessState.ExitCode()
}
