package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const inputTestFile = "testdata/input.txt"

func testOneFile(t *testing.T, srcFileName, tmpFileName, expectedFileName string, offset, limit int64) {
	err := Copy(srcFileName, tmpFileName, offset, limit)
	require.NoError(t, err)
	actual, _ := os.ReadFile(tmpFileName)
	expected, _ := os.ReadFile(expectedFileName)
	require.Equal(t, bytes.Compare(actual, expected), 0)
}

func TestCopy(t *testing.T) {
	t.Run("error cases", func(t *testing.T) {
		err := Copy("", "", 0, 0)
		require.ErrorIs(t, err, ErrEmptyFileName)

		err = Copy("/dev/random", "///", 1, 1)
		if _, ok := err.(*os.PathError); !ok {
			require.Fail(t, "err is not *os.PathError")
		}

		err = Copy("****", "/dev/null", 0, 0)
		if _, ok := err.(*os.PathError); !ok {
			require.Fail(t, "err is not *os.PathError")
		}

		// не сможет записать
		err = Copy("/dev/random", "/dev/full", 0, 1)
		if _, ok := err.(*os.PathError); !ok {
			require.Fail(t, "err is not *os.PathError")
		}

		err = Copy("/dev/random", "/dev/full", -1, 0)
		require.ErrorIs(t, err, ErrInvalidOffset)

		err = Copy("/dev/random", "/dev/full", 0, -1)
		require.ErrorIs(t, err, ErrInvalidLimit)

		err = Copy(".", "/dev/null", 0, 0)
		require.ErrorIs(t, err, ErrUnsupportedFile)

		err = Copy("/dev/random", "/dev/null", 0, 0)
		require.ErrorIs(t, err, ErrInvalidLimitPositive)
	})
	t.Run("success cases", func(t *testing.T) {
		err := Copy("/dev/random", "/dev/null", 0, 100)
		require.NoError(t, err)

		err = Copy("/dev/random", "/dev/null", 100, 100)
		require.NoError(t, err)
	})
	t.Run("test file data", func(t *testing.T) {
		f, _ := os.CreateTemp("", "copy_test")
		defer os.Remove(f.Name())
		f.Close()

		testOneFile(t, inputTestFile, f.Name(), "testdata/out_offset0_limit0.txt", 0, 0)
		testOneFile(t, inputTestFile, f.Name(), "testdata/out_offset0_limit10.txt", 0, 10)
		testOneFile(t, inputTestFile, f.Name(), "testdata/out_offset0_limit1000.txt", 0, 1000)
		testOneFile(t, inputTestFile, f.Name(), "testdata/out_offset0_limit10000.txt", 0, 10000)
		testOneFile(t, inputTestFile, f.Name(), "testdata/out_offset100_limit1000.txt", 100, 1000)
		testOneFile(t, inputTestFile, f.Name(), "testdata/out_offset6000_limit1000.txt", 6000, 1000)
	})

	t.Run("test same file", func(t *testing.T) {
		f, _ := os.CreateTemp("", "copy_test")
		defer os.Remove(f.Name())
		f.Close()
		// копия тестового файла
		testOneFile(t, inputTestFile, f.Name(), "testdata/out_offset0_limit0.txt", 0, 0)
		// сам в себя, полный размер
		testOneFile(t, f.Name(), f.Name(), "testdata/out_offset0_limit0.txt", 0, 0)
		// сам в себя, со смещением и длиной
		testOneFile(t, f.Name(), f.Name(), "testdata/out_offset100_limit1000.txt", 100, 1000)
	})

}
