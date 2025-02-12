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

	exCmd.Env = os.Environ()

	for k, v := range env {
		var value string

		if v.NeedRemove {
			value = ""
		} else {
			value = v.Value
		}

		exCmd.Env = append(exCmd.Env, k+"="+value)
	}

	if err := exCmd.Run(); err != nil {
		log.Println(err.Error())
	}

	return exCmd.ProcessState.ExitCode()
}
