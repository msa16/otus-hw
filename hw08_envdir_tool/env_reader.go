package main

import (
	"bufio"
	"errors"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	result := make(Environment)
	for _, dirEntry := range dirEntries {
		// not directory, and the name must not contain =
		if dirEntry.IsDir() || strings.Contains(dirEntry.Name(), "=") {
			continue
		}
		info, err := dirEntry.Info()
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				// file deleted after creating file list, not error
				continue
			}
			return nil, err
		}
		if info.Size() == 0 {
			// If  the file s is completely empty (0 bytes long), envdir removes an environment variable named s if one exists, without adding a new variable.
			result[dirEntry.Name()] = EnvValue{NeedRemove: true}
		} else {
			// The first line is a value
			file, err := os.Open(path.Join(dir, dirEntry.Name()))
			if err != nil {
				return nil, err
			}
			defer file.Close()

			firstLine := ""
			scanner := bufio.NewScanner(file)
			if scanner.Scan() {
				// Spaces and tabs at the end of t are removed
				firstLine = strings.TrimRight(scanner.Text(), "\t ")
				// Nulls in t are changed	to newlines in the environment variable.
				firstLine = strings.ReplaceAll(firstLine, "\x00", "\n")
			} else {
				if err := scanner.Err(); err != nil {
					return nil, err
				}
			}
			result[dirEntry.Name()] = EnvValue{Value: firstLine, NeedRemove: false}
		}
	}
	return result, nil
}
