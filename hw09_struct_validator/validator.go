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

var allowedConstraints = map[reflect.Kind]map[string]bool{
	reflect.String: {"len": true, "regexp": true, "in": true},
	reflect.Int:    {"min": true, "max": true, "in": true},
}

func checkAllowedConstraints(constraints map[string][]string, kind reflect.Kind) error {
	kindConstraints, ok := allowedConstraints[kind]
	if !ok {
		return nil
	}

	for constraintName := range constraints {
		if _, ok := kindConstraints[constraintName]; !ok {
			return ErrWrongConstraint
		}
	}
	return nil
}

func Validate(v interface{}) error {
	errs := ValidationErrors{}
	val := reflect.ValueOf(v)
	typeInfo := reflect.TypeOf(v)

	for fieldIndex := 0; fieldIndex < val.NumField(); fieldIndex++ {
		fieldInfo := typeInfo.Field(fieldIndex)

		if !fieldInfo.IsExported() {
			continue
		}

		tag := fieldInfo.Tag.Get("validate")
		if tag == "" {
			continue
		}

		checkedValue := val.Field(fieldIndex)

		constraints, err := parseTagFuncs(tag)
		if err != nil {
			errs = append(errs, ValidationError{Field: fieldInfo.Name, Err: err})
		}

		switch checkedValue.Kind() { //nolint:exhaustive
		case reflect.String:
			if err := checkAllowedConstraints(constraints, reflect.String); err != nil {
				return ErrWrongConstraint
			}
			err := validateString(checkedValue, constraints)
			if err != nil {
				errs = append(errs, ValidationError{Field: fieldInfo.Name, Err: err})
			}

		case reflect.Int:
			if err := checkAllowedConstraints(constraints, reflect.Int); err != nil {
				return ErrWrongConstraint
			}
			err := validateInt(checkedValue, constraints)
			if err != nil {
				errs = append(errs, ValidationError{Field: fieldInfo.Name, Err: err})
			}

		case reflect.Slice:
			data := checkedValue.Slice(0, checkedValue.Len())
			validateErrs := validateSlice(constraints, data, fieldInfo.Name)
			if validateErrs == nil {
				continue
			}
			var e ValidationErrors
			if errors.As(validateErrs, &e) {
				errs = append(errs, e...)
			} else {
				return validateErrs
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
) error {
	validErrs := ValidationErrors{}
	if checkedValue.Len() == 0 {
		return nil
	}

	for fieldIndex := 0; fieldIndex < checkedValue.Len(); fieldIndex++ {
		valToCheck := checkedValue.Index(fieldIndex)

		switch valToCheck.Kind() { //nolint:exhaustive
		case reflect.String:
			if err := checkAllowedConstraints(constraints, reflect.String); err != nil {
				return ErrWrongConstraint
			}
			err := validateString(valToCheck, constraints)
			if err != nil {
				validErrs = append(validErrs, ValidationError{Field: fieldName, Err: err})
			}
		case reflect.Int:
			if err := checkAllowedConstraints(constraints, reflect.Int); err != nil {
				return ErrWrongConstraint
			}
			err := validateInt(valToCheck, constraints)
			if err != nil {
				validErrs = append(validErrs, ValidationError{Field: fieldName, Err: err})
			}
		default:
		}
	}
	if len(validErrs) == 0 {
		return nil
	}
	return validErrs
}

var (
	ErrConditionsRepeat = errors.New("conditions are repeated")
	ErrWrongConstraint  = errors.New("wrong constraint")
)

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
