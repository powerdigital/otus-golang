package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
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
		structData  interface{}
		expectedErr error
	}{
		{
			User{
				ID:     "123",
				Name:   "John",
				Age:    25,
				Email:  "john@gmail.com",
				Role:   "stuff",
				Phones: []string{"+123456789"},
				meta:   json.RawMessage("{}"),
			},
			errors.New(strings.Join([]string{
				"Field `ID` contains an error: invalid string length detected",
				"Field `Phones` contains an error: invalid string length detected",
			}, ", ")),
		},
		{
			App{
				Version: "111",
			},
			errors.New("Field `Version` contains an error: invalid string length detected"),
		},
		{
			Response{
				Code: 202,
				Body: "some body",
			},
			errors.New("Field `Code` contains an error: number not in required enum list"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.structData)
			fmtErr := err.Error()
			require.Equal(t, fmtErr, tt.expectedErr.Error())

			_ = tt
		})
	}
}

func TestValidateNoError(t *testing.T) {
	tests := []struct {
		structData interface{}
	}{
		{
			User{
				ID:     "bb829bd2-0c94-4559-ba87-a401291033a8",
				Name:   "John",
				Age:    25,
				Email:  "john@gmail.com",
				Role:   "stuff",
				Phones: []string{"+0123456789"},
				meta:   json.RawMessage("{}"),
			},
		},
		{
			App{
				Version: "11111",
			},
		},
		{
			Response{
				Code: 200,
				Body: "some body",
			},
		},
		{
			Token{
				Header:    nil,
				Payload:   nil,
				Signature: nil,
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.structData)
			require.NoError(t, err)

			_ = tt
		})
	}
}
