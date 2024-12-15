package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testOneExec(t *testing.T, caption string, cmd []string, env Environment, expected int) {
	t.Run(caption, func(t *testing.T) {
		exitCode := RunCmd(cmd, env)
		require.Equal(t, expected, exitCode)
	})

}

func TestRunCmd(t *testing.T) {
	testOneExec(t, "echo", []string{"echo", "hello"}, nil, 0)
	testOneExec(t, "false - code 1", []string{"false"}, nil, 1)
	testOneExec(t, "error starting process - code 111", []string{""}, nil, 111)

	env := Environment{
		"FOO": EnvValue{Value: "foo", NeedRemove: false},
		"BAR": EnvValue{Value: "bar", NeedRemove: false},
	}
	testOneExec(t, "echo with env", []string{"echo", "$FOO $BAR"}, env, 0)
}
