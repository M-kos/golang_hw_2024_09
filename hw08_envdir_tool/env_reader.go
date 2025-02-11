package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
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
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envMap := make(Environment, len(files))
	openedFiles := make([]*os.File, 0)

	defer func() {
		for _, f := range openedFiles {
			f.Close()
		}
	}()

	for _, file := range files {
		if strings.Contains(file.Name(), "=") {
			continue
		}

		f, err := os.Open(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		openedFiles = append(openedFiles, f)

		if s, err := f.Stat(); err != nil || s.Size() == 0 {
			envMap[file.Name()] = EnvValue{
				NeedRemove: true,
			}

			continue
		}

		r := bufio.NewReader(f)
		l, _, err := r.ReadLine()

		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}

		value := clearString(l)

		envMap[file.Name()] = EnvValue{
			Value:      value,
			NeedRemove: value == "",
		}
	}

	return envMap, nil
}

func clearString(line []byte) string {
	s := bytes.TrimRight(line, " \t\n")
	s = bytes.ReplaceAll(s, []byte("\x00"), []byte("\n"))

	return string(s)
}
