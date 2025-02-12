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

		ev, err := fileHandler(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		envMap[file.Name()] = EnvValue{
			Value:      ev.Value,
			NeedRemove: ev.NeedRemove,
		}
	}

	return envMap, nil
}

func clearString(line []byte) string {
	s := bytes.TrimRight(line, " \t\n")
	s = bytes.ReplaceAll(s, []byte("\x00"), []byte("\n"))

	return string(s)
}

func fileHandler(path string) (*EnvValue, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.Size() == 0 {
		return &EnvValue{
			NeedRemove: true,
		}, nil
	}

	r := bufio.NewReader(f)
	l, _, err := r.ReadLine()

	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	value := clearString(l)

	return &EnvValue{
		Value:      value,
		NeedRemove: value == "",
	}, nil
}
