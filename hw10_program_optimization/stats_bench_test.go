package hw10programoptimization

import (
	"archive/zip"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

// go test -bench=. -benchmem .
func getData(b *testing.B) (*zip.ReadCloser, io.ReadCloser) {
	b.Helper()

	r, err := zip.OpenReader("testdata/users.dat.zip")
	require.NoError(b, err)

	require.Equal(b, 1, len(r.File))

	data, err := r.File[0].Open()
	require.NoError(b, err)
	return r, data
}

func BenchmarkGetUsers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r, data := getData(b)
		defer r.Close()
		defer data.Close()

		b.StartTimer()
		getUsers(data)
		b.StopTimer()
	}
}

func BenchmarkCountDomains(b *testing.B) {
	r, data := getData(b)
	defer r.Close()
	defer data.Close()
	u, err := getUsers(data)
	require.NoError(b, err)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		countDomains(u, "biz")
	}
}
