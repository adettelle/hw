package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf("Error in filed: %v, err: %v", ve.Field, ve.Err)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	x := []string{}
	for _, val := range v {
		x = append(x, val.Error())
	}
	slices.Sort(x)
	return strings.Join(x, "; ")
}

func Validate(v interface{}) error {
	errs := ValidationErrors{}
	val := reflect.ValueOf(v)
	typeInfo := reflect.TypeOf(v)

	// Пробегаемся по полям структуры
	for fieldIndex := 0; fieldIndex < val.NumField(); fieldIndex++ {
		fieldInfo := typeInfo.Field(fieldIndex)

		// проверяем, экспортируемое ли поле
		if !fieldInfo.IsExported() {
			continue
		}

		tag := fieldInfo.Tag.Get("validate")
		if tag == "" {
			continue
		}

		// Получаем значение поля структуры, напимер, возраст 19
		checkedValue := val.Field(fieldIndex)

		constraints, err := parseTagFuncs(tag) // ранее tagFuncs
		if err != nil {
			errs = append(errs, ValidationError{Field: fieldInfo.Name, Err: err})
		}

		// проверяем, какое поле int или string (25 - Int)
		// TODO надо ли проверять float???
		switch checkedValue.Kind() {
		case reflect.String:
			err := validateString(checkedValue, constraints) // , validErrs, field
			if err != nil {
				errs = append(errs, ValidationError{Field: fieldInfo.Name, Err: err})
			}

		case reflect.Int:
			err := validateInt(checkedValue, constraints) //  checkRestrictionsInt(tFuncs, int(checkedValue.Int()))
			if err != nil {
				errs = append(errs, ValidationError{Field: fieldInfo.Name, Err: err})
			}

		case reflect.Slice:
			data := checkedValue.Slice(0, checkedValue.Len())
			validateErrs := validateSlice(constraints, data, fieldInfo.Name)
			if len(validateErrs) != 0 {
				errs = append(errs, validateErrs...)
			}
		default:
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateString(checkedValue reflect.Value, constraints map[string][]string) error {
	for k, v := range constraints {
		fmt.Println("K AND V", k, v)
		switch k {
		case "len":
			value, err := strconv.Atoi(v[0])
			if err != nil {
				return err
			}
			if checkedValue.Len() != value {
				return NewErrCheckedValueLen(value)
			}
		case "regexp":
			re, err := regexp.Compile(v[0])
			if err != nil {
				return err
			}
			ok := re.MatchString(checkedValue.String())
			if !ok {
				return NewErrMatch(v[0])
			}
		case "in":
			in := false
			for _, word := range v {
				if checkedValue.String() == word {
					in = true
				}
			}
			if !in {
				return NewErrNotInRange(v)
			}
		}
	}
	return nil
}

func validateInt(checkedValue reflect.Value, constraints map[string][]string) error {
	for k, v := range constraints {
		switch k {
		case "min":
			minValue, err := strconv.Atoi(v[0])
			if err != nil {
				return err
			}
			if int(checkedValue.Int()) < minValue {
				return NewErrValueTooSmall(minValue)
			}
		case "max":
			maxValue, err := strconv.Atoi(v[0])
			if err != nil {
				return err
			}
			if int(checkedValue.Int()) > maxValue {
				return NewErrValueTooBig(maxValue)
			}
		case "in":
			found := false
			for _, val := range v {
				num, err := strconv.Atoi(val)
				if err != nil {
					return err
				}
				if num == int(checkedValue.Int()) {
					found = true
				}
			}
			if !found {
				return NewErrNotInRange(v)
			}
		}
	}
	return nil
}

func validateSlice(constraints map[string][]string,
	checkedValue reflect.Value, fieldName string,
) ValidationErrors {
	validErrs := ValidationErrors{}
	if checkedValue.Len() == 0 {
		for k, v := range constraints {
			num, err := strconv.Atoi(v[0])
			if err != nil {
				return ValidationErrors{ValidationError{Err: err}}
			}
			switch k {
			case "len":
				err = NewErrCheckedValueLen(num)
				ve := ValidationError{Field: fieldName, Err: err}
				validErrs = append(validErrs, ve)
			case "min":
				err := NewErrValueTooSmall(num)
				ve := ValidationError{Field: fieldName, Err: err}
				validErrs = append(validErrs, ve)
			case "max":
				err := NewErrValueTooBig(num)
				ve := ValidationError{Field: fieldName, Err: err}
				validErrs = append(validErrs, ve)
			}
		}
	}

	for fieldIndex := 0; fieldIndex < checkedValue.Len(); fieldIndex++ {
		valToCheck := checkedValue.Index(fieldIndex)

		switch valToCheck.Kind() {
		case reflect.String:
			err := validateString(valToCheck, constraints)
			if err != nil {
				validErrs = append(validErrs, ValidationError{Field: fieldName, Err: err})
			}
		case reflect.Int:

			err := validateInt(valToCheck, constraints) // checkRestrictionsInt(tagFuncs, int(checkedValue.Int()))
			if err != nil {
				validErrs = append(validErrs, ValidationError{Field: fieldName, Err: err})
			}
		}
	}
	if len(validErrs) == 0 {
		return nil
	}
	return validErrs
}

var ErrConditionsRepeat = errors.New("conditions are repeated")

type ErrValueTooBig struct {
	expectedMax int
}

func NewErrValueTooBig(expectedMax int) *ErrValueTooBig {
	return &ErrValueTooBig{expectedMax: expectedMax}
}

func (err *ErrValueTooBig) Error() string {
	return fmt.Sprintf("checkedValue should be less than: %v", err.expectedMax)
}

type ErrValueTooSmall struct {
	expectedMin int
}

func NewErrValueTooSmall(expectedMin int) *ErrValueTooSmall {
	return &ErrValueTooSmall{expectedMin: expectedMin}
}

func (err *ErrValueTooSmall) Error() string {
	return fmt.Sprintf("checkedValue should be grater than: %v", err.expectedMin)
}

type ErrNotInRange struct {
	limit []string
}

func NewErrNotInRange(limit []string) *ErrNotInRange {
	return &ErrNotInRange{limit: limit}
}

func (err *ErrNotInRange) Error() string {
	return fmt.Sprintf("checkedValue should be in range: %v", err.limit)
}

type ErrMatch struct {
	re string
}

func NewErrMatch(re string) *ErrMatch {
	return &ErrMatch{re: re}
}

func (err *ErrMatch) Error() string {
	return fmt.Sprintf("checkedValue should match regular expression: %v", err.re)
}

type ErrCheckedValueLen struct {
	length int
}

func NewErrCheckedValueLen(length int) *ErrCheckedValueLen {
	return &ErrCheckedValueLen{length: length}
}

func (e *ErrCheckedValueLen) Error() string {
	return fmt.Sprintf("wrong length of checkedValue, %v", e.length)
}

func parseTagFuncs(tag string) (map[string][]string, error) {
	funcsMap := map[string][]string{}
	funcs := strings.Split(tag, "|")

	for _, val := range funcs {
		elem := strings.Split(val, ":")
		_, ok := funcsMap[elem[0]]
		if ok {
			return nil, ErrConditionsRepeat
		}
		limits := strings.Split(elem[1], ",")
		funcsMap[elem[0]] = limits
	}
	return funcsMap, nil
}
