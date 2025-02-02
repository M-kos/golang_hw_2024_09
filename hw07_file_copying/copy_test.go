package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	fromPath = "testdata/input.txt"
	toPath = "testdata/test.txt"
)

func TestCopy(t *testing.T) {
	testCases := []struct {
		limit int64
		offset int64
		path string
		name string
	}{
		{limit: 0, offset: 0, path: "testdata/out_offset0_limit0.txt", name: "out_offset0_limit0"},
		{limit: 10, offset: 0, path: "testdata/out_offset0_limit10.txt", name: "out_offset0_limit10"},
		{limit: 1000, offset: 0, path: "testdata/out_offset0_limit1000.txt", name: "out_offset0_limit1000"},
		{limit: 10000, offset: 0, path: "testdata/out_offset0_limit10000.txt", name: "out_offset0_limit10000"},
		{limit: 1000, offset: 100, path: "testdata/out_offset100_limit1000.txt", name: "out_offset100_limit1000"},
		{limit: 1000, offset: 6000, path: "testdata/out_offset6000_limit1000.txt", name: "out_offset6000_limit1000"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := Copy(fromPath, toPath, tt.offset, tt.limit)

			if err != nil {
				t.Error(err.Error())
			}

			require.FileExists(t, toPath)

			testFile, _ := os.Open(toPath)
			masterFile, _ := os.Open(tt.path)

			statTest, _ := testFile.Stat()
			statMaster, _ := masterFile.Stat()

			require.Equal(t, statMaster.Size(), statTest.Size())

			os.Remove(toPath)
		})
	}
}

func TestCopyError(t *testing.T) {
	var limit int64 = 0
	var offset int64 = 7000

	t.Run("Check offset exceeds error", func(t *testing.T) {
			err := Copy(fromPath, toPath, offset, limit)

			require.ErrorIs(t, ErrOffsetExceedsFileSize, err)
		})
}
