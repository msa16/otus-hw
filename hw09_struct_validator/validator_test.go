package hw09structvalidator

import (
	"encoding/json"
	"errors"
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

	Bad1 struct {
		Phone string `validate:"len:X"`
	}
	Bad2 struct {
		Age int `validate:"min:X"`
	}
	Bad3 struct {
		Age int `validate:"max:X"`
	}
	Bad4 struct {
		Age int `validate:"in:10,X,30"`
	}
	Bad5 struct {
		Age int `validate:"max"`
	}
	Bad6 struct {
		Phone string `validate:"regexp:Hello(?|!)"`
	}
	Bad7 struct {
		Age int `validate:"zero:0"`
	}
	IntArray struct {
		Data []int `validate:"min:777|max:999"`
	}
	Staff struct {
		User `validate:"nested"`
		Unit string `validate:"in:бухгалтерия,транспортный цех"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          User{ID: "12345678-90ab-1234-5678-0123456789ab", Name: "Somebody", Age: 25, Email: "some@ya.ru", Role: "admin", Phones: []string{"79169999999", "79169999998"}},
			expectedErr: nil,
		},
		{
			in:          0,
			expectedErr: ErrArgumentNotStructure,
		},
		{
			in:          Bad1{Phone: "+79169999999"},
			expectedErr: ErrInvalidLenTag,
		},
		{
			in:          Bad2{Age: 10},
			expectedErr: ErrInvalidMinTag,
		},
		{
			in:          Bad3{Age: 10},
			expectedErr: ErrInvalidMaxTag,
		},
		{
			in:          Bad4{Age: 1},
			expectedErr: ErrInvalidInTag,
		},
		{
			in:          Bad5{Age: 1},
			expectedErr: ErrInvalidTag,
		},
		{
			in:          Bad6{Phone: "1"},
			expectedErr: ErrInvalidRegExp,
		},
		{
			in:          App{Version: "1234"},
			expectedErr: ValidationErrors{{Field: "Version", Err: ErrWrongLen}},
		},
		{
			in:          App{Version: "АБВГД"},
			expectedErr: nil,
		},
		{
			in:          User{ID: "12345678-90ab-1234-5678-0123456789ab", Name: "Somebody", Age: 25, Email: "some#ya.ru", Role: "admin", Phones: []string{"79169999999", "79169999998"}},
			expectedErr: ValidationErrors{{Field: "Email", Err: ErrRegExpNotMatch}},
		},
		{
			in:          User{ID: "12345678-90ab-1234-5678-0123456789ab", Name: "Somebody", Age: 25, Email: "some@ya.ru", Role: "user", Phones: []string{"79169999999", "79169999998"}},
			expectedErr: ValidationErrors{{Field: "Role", Err: ErrNotInSet}},
		},
		{
			in:          User{ID: "12345678-90ab-1234-5678-0123456789ab", Name: "Somebody", Age: 10, Email: "some@ya.ru", Role: "admin", Phones: []string{"79169999999", "79169999998"}},
			expectedErr: ValidationErrors{{Field: "Age", Err: ErrMinFailed}},
		},
		{
			in:          User{ID: "12345678-90ab-1234-5678-0123456789ab", Name: "Somebody", Age: 70, Email: "some@ya.ru", Role: "admin", Phones: []string{"79169999999", "79169999998"}},
			expectedErr: ValidationErrors{{Field: "Age", Err: ErrMaxFailed}},
		},
		{
			in:          Bad7{Age: 10},
			expectedErr: ErrUnknownValidator,
		},
		{
			in:          Token{Header: []byte("header"), Payload: []byte("payload"), Signature: []byte("signature")},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 333,
			},
			expectedErr: ValidationErrors{{Field: "Code", Err: ErrNotInSet}},
		},
		{
			in:          User{ID: "12345678-90ab-1234-5678-0123456789ab", Name: "Somebody", Age: 25, Email: "some@ya.ru", Role: "admin", Phones: []string{"1", "12"}},
			expectedErr: ValidationErrors{{Field: "Phones", Err: ErrWrongLen}, {Field: "Phones", Err: ErrWrongLen}},
		},
		{
			in: User{ID: "1", Name: "Somebody", Age: 70, Email: "some#ya.ru", Role: "ZZZ", Phones: []string{"+79169999999", "+79169998"}},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrWrongLen},
				{Field: "Age", Err: ErrMaxFailed},
				{Field: "Email", Err: ErrRegExpNotMatch},
				{Field: "Role", Err: ErrNotInSet},
				{Field: "Phones", Err: ErrWrongLen},
				{Field: "Phones", Err: ErrWrongLen}},
		},
		{in: IntArray{Data: []int{777, 999}}, expectedErr: nil},
		{in: IntArray{Data: []int{77, 9999}}, expectedErr: ValidationErrors{{Field: "Data", Err: ErrMinFailed}, {Field: "Data", Err: ErrMaxFailed}}},
		{in: Staff{
			User: User{ID: "0", Name: "Somebody", Age: 75, Email: "some@ya.ru", Role: "admin", Phones: []string{"79169999999", "79169999998"}},
			Unit: "склад"},
			expectedErr: ValidationErrors{{Field: "ID", Err: ErrWrongLen}, {Field: "Age", Err: ErrMaxFailed}, {Field: "Unit", Err: ErrNotInSet}}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			validationErr := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, validationErr)
			} else {
				// массив ValidationErrors
				var expectedErrors ValidationErrors
				if errors.As(tt.expectedErr, &expectedErrors) {
					var actualErrors ValidationErrors
					if errors.As(validationErr, &actualErrors) {
						require.Equal(t, len(expectedErrors), len(actualErrors))
						for i, actualErr := range actualErrors {
							expectedErr := expectedErrors[i]
							require.Equal(t, actualErr.Field, expectedErr.Field)
							require.ErrorIs(t, actualErr.Err, expectedErr.Err)
						}
					} else {
						require.Fail(t, fmt.Sprintf("expected ValidationErrors, but returned %v", validationErr))
					}
				} else {
					// одна программная ошибка
					require.ErrorIs(t, validationErr, tt.expectedErr)
				}
			}
			_ = tt
		})
	}
}
