package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
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
		AvRating []int `validate:"min:3|max:5"`
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

	Cat struct {
		Name string `validate:"len:5|min:3"`
		Age  int
	}

	Dog struct {
		Name    string
		Rewards []int `validate:"max:10|regexp:\\d+"`
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

func TestNegative(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: Cat{
				Name: "Kitty",
				Age:  5,
			},
			expectedErr: ErrWrongConstraint,
		},
		{
			in: Dog{
				Name:    "Snoopy",
				Rewards: []int{10, 9, 8, 7},
			},
			expectedErr: ErrWrongConstraint,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			require.Error(t, err)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestArrangeTagFuncs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		elemType reflect.Type
		expected []FieldValidator
	}{
		{
			name:     "check one condition",
			input:    "len:36",
			elemType: reflect.TypeOf(""),
			expected: []FieldValidator{&LenValidator{Len: 36}},
		},
		{
			name:     "check two conditions",
			input:    "min:18|max:50",
			elemType: reflect.TypeOf(int(1)),
			expected: []FieldValidator{
				&MinValidator{Min: 18},
				&MaxValidator{Max: 50},
			},
		},
		{
			name:     "condition with backslash",
			input:    "regexp:^\\w+@\\w+\\.\\w+$",
			elemType: reflect.TypeOf(""),
			expected: []FieldValidator{&RegexpValidator{Re: regexp.MustCompile(`^\w+@\w+\.\w+$`)}},
		},
		{
			name:     "condition with slice in value",
			input:    "in:admin,stuff",
			elemType: reflect.TypeOf(""),
			expected: []FieldValidator{&InStrValidator{Elems: []string{"admin", "stuff"}}},
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			result, err := parseTagFuncs(tCase.input, tCase.elemType)
			require.NoError(t, err)
			require.Equal(t, tCase.expected, result)
		})
	}
}
