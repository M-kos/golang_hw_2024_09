package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileNOtExist = errors.New("file not exist")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	inFile, err := os.Open(fromPath)

	defer inFile.Close()

	stat, err := inFile.Stat()

	if err != nil {
		return err
	}

	if stat.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNOtExist
		}

		return ErrUnsupportedFile
	}

	_, err = inFile.Seek(offset, io.SeekStart) 

	if err != nil {
		return err
	}

	if limit == 0 {
		limit = stat.Size()
	}

	fmt.Println("Limit: ", limit)

	err = writeToFile(inFile, toPath, offset, limit)

	return nil
}

func writeToFile(inFile io.Reader, toPath string, offset, limit int64) error {
	outFile, err := os.Create(toPath)

	if err != nil {
		return err
	}

	bar := pb.Start64(limit)
	barWriter := bar.NewProxyWriter(outFile)

	_, err = io.CopyN(barWriter, inFile, limit)
	bar.Finish()

	if err != nil {
		if err == io.EOF {
			return nil
		}

		return err
	}

	return nil
}