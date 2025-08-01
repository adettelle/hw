package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return -1
	}

	args := []string{}
	if len(cmd) > 1 {
		args = cmd[1:]
	}

	for name, envVal := range env {
		if envVal.NeedRemove {
			err := os.Unsetenv(name)
			if err != nil {
				return -1
			}
		}

		err := os.Setenv(name, envVal.Value)
		if err != nil {
			return -1
		}
	}

	command := exec.Command(cmd[0], args...)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			return ee.ExitCode()
		}
		return -1
	}
	return 0
}
