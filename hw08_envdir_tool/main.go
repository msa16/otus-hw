package main

import (
	"fmt"
	"os"
)

func main() {
	// https://www.unix.com/man-page/debian/8/envdir/
	args := os.Args
	if len(args) < 3 {
		fmt.Println(ERROR_USAGE)
		os.Exit(ERROR_FATAL_CODE)
	}

	envs, err := ReadDir(args[1])
	if err != nil {
		fmt.Printf(ERROR_FATAL, err)
		os.Exit(ERROR_FATAL_CODE)
	}
	if exitCode := RunCmd(args[2:], envs); exitCode != 0 {
		os.Exit(exitCode)
	}
}
