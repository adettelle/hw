package main

import (
	"os"
	"path"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	inPath    = "./testdata/in/input.txt"
	chunkSize = 3
)

// func TestCopy(t *testing.T) {
// 	// Place your code here.
// }

func TestCopyEntireFileZeroOffset(t *testing.T) {
	progressCh := make(chan int, 5)
	defer close(progressCh)

	completionCh := make(chan any)
	defer close(completionCh)

	wg := sync.WaitGroup{}

	tCase := struct {
		name   string
		from   string
		to     string
		limit  int64
		offset int64
	}{
		name:   "copy_entire_file",
		from:   inPath,
		limit:  0,
		offset: 0,
	}

	toDir := createAndCleanOutDir(t)
	toFile := path.Join(toDir, tCase.name)

	progress(progressCh, completionCh, &wg)

	err := Copy(tCase.from, toFile, tCase.offset, tCase.limit, progressCh, completionCh, chunkSize) // tCase.to
	require.NoError(t, err)

	incomingFileInfo, err := os.Stat(tCase.from)
	require.NoError(t, err)
	outcomigFileInfo, err := os.Stat(toFile)
	require.NoError(t, err)

	require.Equal(t, outcomigFileInfo.Size(), incomingFileInfo.Size())

	wg.Wait()
}

func TestCopyWithOffset(t *testing.T) {
	tests := []struct {
		name       string
		from       string
		to         string
		limit      int64
		offset     int64
		chunkSize  int
		expected   string
		shouldFail bool
	}{
		{
			name:      "copy_with_limit_and_offset",
			from:      inPath,
			limit:     10,
			offset:    5,
			chunkSize: chunkSize,
			expected:  "67890\nabcd",
		},
		{
			name:      "copy_with_offset",
			from:      inPath,
			limit:     0,
			offset:    5,
			chunkSize: chunkSize,
			expected:  "67890\nabcdefghijklmnopqrstuvwxyz",
		},
		{
			name:      "copy_with_enormous_limit",
			from:      inPath,
			limit:     1000,
			offset:    6,
			chunkSize: chunkSize,
			expected:  "7890\nabcdefghijklmnopqrstuvwxyz",
		},
		{
			name:      "copy_with_enormous_chunkSize",
			from:      inPath,
			limit:     0,
			offset:    0,
			chunkSize: 512 * 1024,
			expected:  "1234567890\nabcdefghijklmnopqrstuvwxyz",
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			progressCh := make(chan int, 5)
			defer close(progressCh)

			completionCh := make(chan any, 1)
			defer close(completionCh)

			wg := sync.WaitGroup{}

			toDir := createAndCleanOutDir(t)
			toFile := path.Join(toDir, tCase.name)

			progress(progressCh, completionCh, &wg)

			err := Copy(tCase.from, toFile, tCase.offset, tCase.limit, progressCh, completionCh, tCase.chunkSize) // tCase.to
			if tCase.shouldFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				outFileContent, err := os.ReadFile(toFile)
				require.NoError(t, err)

				require.Equal(t, tCase.expected, string(outFileContent))
			}
			wg.Wait()
		})
	}
}

type test struct {
	from   string
	to     string
	limit  int64
	offset int64
}

func TestCopyToNonexistentDir(t *testing.T) {
	progressCh := make(chan int, 5)
	defer close(progressCh)

	completionCh := make(chan any)
	defer close(completionCh)

	wg := sync.WaitGroup{}

	tCase := test{
		from:   inPath,
		to:     "../non-existentdir/output.txt",
		limit:  0,
		offset: 0,
	}

	_, err := os.Stat(tCase.to)
	require.True(t, os.IsNotExist(err))

	progress(progressCh, completionCh, &wg)

	err = Copy(tCase.from, tCase.to, tCase.offset, tCase.limit, progressCh, completionCh, chunkSize)
	require.Error(t, err, ErrDestinationDirDoesnotExist)

	wg.Wait()
}

func createAndCleanOutDir(t *testing.T) string {
	t.Helper()

	dir, err := os.MkdirTemp("", "testout*")
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})

	return dir
}

func progress(progressCh <-chan int, completionCh <-chan any, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-progressCh:
			case <-completionCh:
				return
			}
		}
	}()
}
