package main

import (
	"os"
	"os/exec"
)

const (
	// envdir exits 111 if it has trouble reading d, if it runs out of memory for environment variables,
	// or if it cannot run child.  Otherwise its exit code is the same as that of child.
	ERROR_FATAL_CODE = 111
	ERROR_FATAL      = "envdir: fatal: %s\n"
	ERROR_USAGE      = "envdir: usage: envdir dir child"
)

func modifyEnv(env Environment) {
	for k, v := range env {
		if v.NeedRemove {
			os.Unsetenv(k)
		} else {
			os.Setenv(k, v.Value)
		}
	}
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	modifyEnv(env)

	if err := command.Start(); err != nil {
		return ERROR_FATAL_CODE
	}

	if err := command.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			return exiterr.ExitCode()
		} else {
			return ERROR_FATAL_CODE
		}
	}
	return 0
}
