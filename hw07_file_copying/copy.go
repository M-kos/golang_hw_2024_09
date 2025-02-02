package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3" //nolint:all
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileNOtExist          = errors.New("file not exist")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	inFile, err := os.Open(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNOtExist
		}

		return ErrUnsupportedFile
	}

	defer inFile.Close()

	stat, err := inFile.Stat()
	if err != nil {
		return err
	}

	if stat.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	_, err = inFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	if limit == 0 {
		limit = stat.Size()
	}

	err = writeToFile(inFile, toPath, limit)
	if err != nil {
		return err
	}

	return nil
}

func writeToFile(inFile io.Reader, toPath string, limit int64) error {
	outFile, err := os.Create(toPath)
	if err != nil {
		return err
	}

	bar := pb.Start64(limit)
	barWriter := bar.NewProxyWriter(outFile)

	_, err = io.CopyN(barWriter, inFile, limit)
	bar.Finish()

	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	return nil
}
