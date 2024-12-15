package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	TEST_DIR             = "dir1"
	TEST_FILE            = "file1"
	TEST_FILE_WRONG_NAME = "file=1"
	TEST_TEXT            = "text"
)

func checkEnv(t *testing.T, env EnvValue, value string, needRemove bool) {
	require.Equal(t, value, env.Value)
	require.Equal(t, needRemove, env.NeedRemove)

}

func TestReadDir(t *testing.T) {
	t.Run("read dir with files", func(t *testing.T) {
		envs, err := ReadDir("testdata/env")
		require.NoError(t, err)
		require.Len(t, envs, 5)

		checkEnv(t, envs["BAR"], "bar", false)
		checkEnv(t, envs["EMPTY"], "", false)
		checkEnv(t, envs["FOO"], "   foo\nwith new line", false)
		checkEnv(t, envs["HELLO"], "\"hello\"", false)
		checkEnv(t, envs["UNSET"], "", true)
	})
}

func TestReadDirExt(t *testing.T) {
	var tempDir string = ""
	t.Cleanup(func() {
		if tempDir != "" {
			err := os.RemoveAll(tempDir)
			require.NoError(t, err)
		}
	})
	tempDir, err := os.MkdirTemp("", "env")
	require.NoError(t, err)

	t.Run("read empty dir", func(t *testing.T) {
		envs, err := ReadDir(tempDir)
		require.NoError(t, err)
		require.Len(t, envs, 0)
	})

	t.Run("read not existing dir", func(t *testing.T) {
		envs, err := ReadDir(path.Join(tempDir, TEST_DIR))
		require.Error(t, err)
		require.Nil(t, envs)
	})

	err = os.Mkdir(path.Join(tempDir, TEST_DIR), 0755)
	require.NoError(t, err)
	err = os.WriteFile(path.Join(tempDir, TEST_FILE), []byte(TEST_TEXT), 0644)
	require.NoError(t, err)
	err = os.WriteFile(path.Join(tempDir, TEST_FILE_WRONG_NAME), []byte(TEST_TEXT), 0644)
	require.NoError(t, err)
	err = os.WriteFile(path.Join(tempDir, TEST_DIR, TEST_FILE), []byte(TEST_TEXT), 0644)
	require.NoError(t, err)

	t.Run("read dir with subdir and file with wrong name", func(t *testing.T) {
		envs, err := ReadDir(tempDir)
		require.NoError(t, err)
		require.Len(t, envs, 1)
		checkEnv(t, envs[TEST_FILE], TEST_TEXT, false)
	})
}
