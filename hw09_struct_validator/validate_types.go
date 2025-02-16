package hw09structvalidator

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Validator interface {
	Valid(value string) error
}

type In struct {
	Values []string
}

func (v In) Valid(value string) error {
	var errValues strings.Builder

	for i, val := range v.Values {
		if val == value {
			return nil
		}

		errValues.WriteString(val)

		if i != len(v.Values)-1 {
			errValues.WriteString(", ")
		}
	}

	errStr := fmt.Sprintf("not in %v", errValues.String())

	return errors.New(errStr)
}

type Len struct {
	Length int
}

func (v Len) Valid(value string) error {
	if len(value) != v.Length {
		return fmt.Errorf("len must be %v", v.Length)
	}

	return nil
}

type Min struct {
	Min int
}

func (v Min) Valid(value string) error {
	convertedValue, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	if convertedValue < v.Min {
		return fmt.Errorf("must be more than %v", v.Min)
	}

	return nil
}

type Max struct {
	Max int
}

func (v Max) Valid(value string) error {
	convertedValue, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	if convertedValue > v.Max {
		return fmt.Errorf("must be less than %v", v.Max)
	}

	return nil
}

type Regexp struct {
	Regexp string
}

func (v Regexp) Valid(value string) error {
	matched, err := regexp.MatchString(v.Regexp, value)
	if err != nil {
		return err
	}

	if !matched {
		return fmt.Errorf("must match %v", v.Regexp)
	}

	return nil
}
