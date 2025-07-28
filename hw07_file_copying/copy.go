package main

import (
	"errors"
	"io"
	"math"
	"os"
	"path"
)

var (
	ErrUnsupportedFile            = errors.New("unsupported file")
	ErrOffsetExceedsFileSize      = errors.New("offset exceeds file size")
	ErrSamePaths                  = errors.New("paths from and to are the same")
	ErrNoSuchFile                 = errors.New("there is no source file")
	ErrIrregularFile              = errors.New("file with unknown size")
	ErrDestinationDirDoesnotExist = errors.New("destination directory does not exist")
)

func Copy(fromPath, toPath string, offset, limit int64,
	progressCh chan<- int, completionCh chan<- any, chunkSize int,
) error {
	defer func() {
		completionCh <- struct{}{}
	}()

	if fromPath == toPath {
		return ErrSamePaths
	}

	sourceFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Move to specific position
	_, err = sourceFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	sourceFi, err := os.Stat(fromPath) // testing whether the file that will be copied exists
	if err != nil {
		return ErrNoSuchFile
	}

	if !sourceFi.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	_, err = os.Stat(path.Dir(toPath))
	if err != nil {
		if os.IsNotExist(err) {
			return ErrDestinationDirDoesnotExist
		}
		return err
	}

	// offset больше, чем размер файла - невалидная ситуация
	if offset > sourceFi.Size() {
		return ErrOffsetExceedsFileSize
	}

	destination, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer destination.Close()

	sourceFileSize := sourceFi.Size()
	if limit == 0 {
		limit = sourceFileSize - offset
	}
	destinationFileSize := min(sourceFileSize-offset, limit)

	numChunks := int(math.Ceil(float64(destinationFileSize) / float64(chunkSize)))

	var lastChunkSize int
	if destinationFileSize != 0 {
		lastChunkSize = int(destinationFileSize) % chunkSize
	}

	for i := range numChunks {
		if i == numChunks-1 {
			chunkSize = lastChunkSize
		}
		_, err := io.CopyN(destination, sourceFile, int64(chunkSize))

		percentage := 100 * (i + 1) / numChunks

		progressCh <- percentage

		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
	}

	return nil
}
