package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	Student struct {
		Name     string
		Age      int   `validate:"max:80"`
		AvRating []int `validate:"len:6|min:3|max:5"`
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

func TestValidatePositive(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Ane",
				Email:  "ane@gmail.com",
				Age:    19,
				Role:   "stuff",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Err: nil,
				},
			},
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Alex",
				Age:    18,
				Email:  "alex@gmail.com",
				Role:   "admin",
				Phones: []string{"12345678901", "12345678901"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Err: nil,
				},
			},
		},
		{
			in: Student{
				Age:      10,
				AvRating: []int{5, 5, 3, 3, 4, 4},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Err: nil,
				},
			},
		},
		{
			in: Student{
				Name:     "Alex",
				AvRating: []int{5, 5, 5, 5, 4, 4},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Err: nil,
				},
			},
		},
		{
			in: App{
				Version: "abcde",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Err: nil,
				},
			},
		},
		{
			in: Token{
				Header:    []byte(""),
				Payload:   []byte(""),
				Signature: []byte(""),
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Err: nil,
				},
			},
		},
		{
			in: Response{
				Code: 404,
				Body: "",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Err: nil,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.NoError(t, err)
		})
	}
}

func TestValidateNegative(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Age:    20,
				Email:  "mary@gmail.com",
				Phones: []string{"1234", "12345678901"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Role",
					Err:   NewErrNotInRange([]string{"admin", "stuff"}),
				},
				ValidationError{
					Field: "Phones",
					Err:   NewErrCheckedValueLen(11),
				},
			},
		},
		{
			in: Student{
				Age: 91,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   NewErrValueTooBig(80),
				},
				ValidationError{
					Field: "AvRating",
					Err:   NewErrCheckedValueLen(6),
				},
				ValidationError{
					Field: "AvRating",
					Err:   NewErrValueTooSmall(3),
				},
				ValidationError{
					Field: "AvRating",
					Err:   NewErrValueTooBig(5),
				},
			},
		},
		{
			in: User{
				Name:   "Ane",
				Email:  "anegmail.com",
				Age:    19,
				Phones: []string{"1234", "12345678901"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   NewErrCheckedValueLen(36),
				},
				ValidationError{
					Field: "Email",
					Err:   NewErrMatch("^\\w+@\\w+\\.\\w+$"),
				},
				ValidationError{
					Field: "Role",
					Err:   NewErrNotInRange([]string{"admin", "stuff"}),
				},
				ValidationError{
					Field: "Phones",
					Err:   NewErrCheckedValueLen(11),
				},
			},
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Alex",
				Age:    17,
				Email:  "alex@gmail.com",
				Role:   "admin",
				Phones: []string{"12345678901", "12345"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   NewErrValueTooSmall(18),
				},
				ValidationError{
					Field: "Phones",
					Err:   NewErrCheckedValueLen(11),
				},
			},
		},
		{
			in: App{
				Version: "abcdefg",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   NewErrCheckedValueLen(5),
				},
			},
		},
		{
			in: Response{
				Code: 400,
				Body: "",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   NewErrNotInRange([]string{"200", "404", "500"}),
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Error(t, err)
			require.Equal(t, tt.expectedErr.Error(), err.Error())
		})
	}
}

func TestArrangeTagFuncs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string][]string
	}{
		{
			name:     "check one condition",
			input:    "len:36",
			expected: map[string][]string{"len": {"36"}},
		},
		{
			name:     "check two conditions",
			input:    "min:18|max:50",
			expected: map[string][]string{"min": {"18"}, "max": {"50"}},
		},
		{
			name:     "condition with backslash",
			input:    "regexp:^\\w+@\\w+\\.\\w+$",
			expected: map[string][]string{"regexp": {"^\\w+@\\w+\\.\\w+$"}},
		},
		{
			name:     "condition with slice in value",
			input:    "in:admin,stuff",
			expected: map[string][]string{"in": {"admin", "stuff"}},
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			result, err := parseTagFuncs(tCase.input)
			require.NoError(t, err)
			require.Equal(t, tCase.expected, result)
		})
	}
}
