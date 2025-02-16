package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require" //nolint:all
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				meta:   json.RawMessage(`{"key": "value"}`),
				ID:     "123456789012345678901234567890123456",
				Name:   "John",
				Age:    25,
				Email:  "john@example.com",
				Role:   "admin",
				Phones: []string{"89999999999", "89999999998"},
			},
			expectedErr: nil,
		},
		{
			in:          App{Version: "1.0.0"},
			expectedErr: nil,
		},
		{
			in:          Token{},
			expectedErr: nil,
		},
		{
			in:          Response{Code: 200, Body: "OK"},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "123",
				Age:    5,
				Email:  "qwe",
				Role:   "stuff",
				Phones: []string{"+7"},
			},
			expectedErr: errors.New("ID: len must be 36; Age: must be more than 18; Email: must match ^\\w+@\\w+\\.\\w+$; Phones: len must be 11; "), //nolint:all
		},
		{
			in: Response{
				Code: 100,
			},
			expectedErr: errors.New("Code: not in 200, 404, 500; "),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, tt.expectedErr.Error(), err.Error())
			}
		})
	}
}
