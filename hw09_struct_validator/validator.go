package hw09structvalidator

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var errors strings.Builder
	for _, err := range v {
		errors.WriteString(err.Field)
		errors.WriteString(": ")
		errors.WriteString(err.Err.Error())
		errors.WriteString("; ")
	}

	return errors.String()
}

const (
	TagValidate    = "validate"
	SplitTagValues = "|"
)

var (
	ErrNotAStruct = errors.New("not a struct")
	ErrUnknownTag = errors.New("unknown tag")
)

func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	errs := ValidationErrors{}

	if rv.Kind() != reflect.Struct {
		return ErrNotAStruct
	}

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Type().Field(i)
		tag := field.Tag.Get(TagValidate)

		if tag == "" {
			continue
		}

		if err := validateField(rv.Field(i), tag); err != nil {
			errs = append(errs, ValidationError{Field: field.Name, Err: err})
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func validateField(field reflect.Value, tag string) error {
	tagValues := strings.Split(tag, SplitTagValues)
	var err error

	for _, tagValue := range tagValues {
		name, tagValue := parseTag(tagValue)
		validator, err := getValidator(name, tagValue)
		if err != nil {
			return err
		}

		if field.Kind() == reflect.String {
			err = validator.Valid(field.String())
			return err
		}

		if field.Kind() == reflect.Int {
			val := int(field.Int())
			err = validator.Valid(strconv.Itoa(val))
			return err
		}

		if field.Kind() == reflect.Slice {
			err = sliceValidate(field, validator)
			return err
		}
	}

	return err
}

func getValidator(name, tagValue string) (Validator, error) {
	switch name {
	case "in":
		values := strings.Split(tagValue, ",")
		return In{
			Values: values,
		}, nil
	case "len":
		value, err := strconv.Atoi(tagValue)
		if err != nil {
			return nil, err
		}

		return Len{
			Length: value,
		}, nil
	case "min":
		value, err := strconv.Atoi(tagValue)
		if err != nil {
			return nil, err
		}

		return Min{
			Min: value,
		}, nil
	case "max":
		value, err := strconv.Atoi(tagValue)
		if err != nil {
			return nil, err
		}

		return Max{
			Max: value,
		}, nil
	case "regexp":
		return Regexp{
			Regexp: tagValue,
		}, nil
	default:
		return nil, ErrUnknownTag
	}
}

func parseTag(tagValue string) (string, string) {
	parts := strings.Split(tagValue, ":")

	return parts[0], parts[1]
}

func sliceValidate(field reflect.Value, validator Validator) error {
	for i := 0; i < field.Len(); i++ {
		rv := field.Index(i)

		if err := validator.Valid(rv.String()); err != nil {
			return err
		}
	}

	return nil
}
