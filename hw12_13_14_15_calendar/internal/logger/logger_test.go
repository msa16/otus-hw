package logger

import (
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	checkLevel(t, "DEBUG")
	checkLevel(t, "INFO")
	checkLevel(t, "WARN")
	checkLevel(t, "ERROR")
}

func checkLevel(t *testing.T, level string) {
	t.Helper()
	t.Run("level "+level, func(t *testing.T) {
		f, _ := os.CreateTemp("", "logger_test_"+level+".log")
		defer require.NoError(t, os.Remove(f.Name()))
		f.Close()

		logg := New(level, f.Name())
		if _, err := os.Stat(f.Name()); errors.Is(err, fs.ErrNotExist) {
			t.Errorf("log file does not exist")
		}
		logg.Debug("test debug")
		logg.Info("test info")
		logg.Warn("test warn")
		logg.Error("test error")
		require.NoError(t, logg.Close())
	})
}
