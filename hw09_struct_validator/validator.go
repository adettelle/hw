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

		constraints, err := parseTagFuncs(tag, fieldInfo.Type)
		if err != nil {
			return err // parse error
		}
		if len(constraints) == 0 {
			continue
		}

		if checkedValue.Kind() == reflect.Slice {
			for fieldIndex := 0; fieldIndex < checkedValue.Len(); fieldIndex++ {
				valToCheck := checkedValue.Index(fieldIndex)
				for _, constraint := range constraints {
					err := constraint.validateValue(valToCheck)
					if err != nil {
						errs = append(errs, ValidationError{Field: fieldInfo.Name, Err: err})
					}
				}
			}
		} else {
			for _, constraint := range constraints {
				err := constraint.validateValue(checkedValue)
				if err != nil {
					errs = append(errs, ValidationError{Field: fieldInfo.Name, Err: err})
				}
			}
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
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

type ErrNotInIntRange struct {
	limit []int
}

func NewErrNotInIntRange(limit []int) *ErrNotInIntRange {
	return &ErrNotInIntRange{limit: limit}
}

func (err *ErrNotInIntRange) Error() string {
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

func parseTagFuncs(tag string, rt reflect.Type) ([]FieldValidator, error) {
	elems := strings.Split(tag, "|") // elems is []string{"min:1", "max:2"}

	fvs := []FieldValidator{}
	kind := rt.Kind()
	if kind == reflect.Slice {
		kind = rt.Elem().Kind()
	}
	for _, elem := range elems { // elem is min:1 or max:2
		var v FieldValidator
		var err error

		funcAndVal := strings.Split(elem, ":")
		switch funcAndVal[0] {
		case "max":
			v, err = NewMaxValidator(funcAndVal[1])
		case "min":
			v, err = NewMinValidator(funcAndVal[1])
		case "in":
			switch kind { //nolint:exhaustive
			case reflect.Int:
				v, err = NewInNumValidator(funcAndVal[1])
			case reflect.String:
				v, err = NewInStrValidator(funcAndVal[1])
			default:
				return nil, ErrWrongConstraint
			}
		case "len":
			v, err = NewLenValidator(funcAndVal[1])
		case "regexp":
			v, err = NewRegexpValidator(funcAndVal[1])
		default:
			return nil, ErrWrongConstraint
		}
		if err != nil {
			return nil, err
		}
		if !v.acsepts(kind) {
			return nil, ErrWrongConstraint
		}
		fvs = append(fvs, v)
	}
	return fvs, nil
}

type FieldValidator interface {
	acsepts(reflect.Kind) bool
	validateValue(reflect.Value) error
}

type MaxValidator struct {
	Max int
}

func NewMaxValidator(s string) (*MaxValidator, error) {
	maxVal, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	return &MaxValidator{Max: maxVal}, nil
}

func (mv *MaxValidator) acsepts(kind reflect.Kind) bool {
	return kind == reflect.Int
}

func (mv *MaxValidator) validateValue(item reflect.Value) error {
	if int(item.Int()) > mv.Max {
		return NewErrValueTooBig(mv.Max)
	}
	return nil
}

type MinValidator struct {
	Min int
}

func NewMinValidator(s string) (*MinValidator, error) {
	minVal, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	return &MinValidator{Min: minVal}, nil
}

func (mv *MinValidator) acsepts(kind reflect.Kind) bool {
	return kind == reflect.Int
}

func (mv *MinValidator) validateValue(item reflect.Value) error {
	if int(item.Int()) < mv.Min {
		return NewErrValueTooSmall(mv.Min)
	}
	return nil
}

type InNumValidator struct {
	Nums []int
}

func NewInNumValidator(s string) (*InNumValidator, error) {
	elems := strings.Split(s, ",")
	nums := []int{}
	for _, elem := range elems {
		n, err := strconv.Atoi(elem)
		if err != nil {
			return nil, err
		}
		nums = append(nums, n)
	}
	return &InNumValidator{Nums: nums}, nil
}

func (mv *InNumValidator) acsepts(kind reflect.Kind) bool {
	return kind == reflect.Int
}

func (mv *InNumValidator) validateValue(item reflect.Value) error {
	hasElement := slices.Index(mv.Nums, int(item.Int())) > -1
	if !hasElement {
		return NewErrNotInIntRange(mv.Nums)
	}
	return nil
}

type InStrValidator struct {
	Elems []string
}

func NewInStrValidator(s string) (*InStrValidator, error) {
	elems := strings.Split(s, ",")
	return &InStrValidator{Elems: elems}, nil
}

func (mv *InStrValidator) acsepts(kind reflect.Kind) bool {
	return kind == reflect.String
}

func (mv *InStrValidator) validateValue(item reflect.Value) error {
	hasElement := slices.Index(mv.Elems, item.String()) > -1
	if !hasElement {
		return NewErrNotInRange(mv.Elems)
	}
	return nil
}

type LenValidator struct {
	Len int
}

func NewLenValidator(s string) (*LenValidator, error) {
	lenVal, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	return &LenValidator{Len: lenVal}, nil
}

func (lv *LenValidator) acsepts(kind reflect.Kind) bool {
	return kind == reflect.String
}

func (lv *LenValidator) validateValue(item reflect.Value) error {
	if len(item.String()) != lv.Len {
		return NewErrCheckedValueLen(lv.Len)
	}
	return nil
}

type RegexpValidator struct {
	Re *regexp.Regexp
}

func NewRegexpValidator(s string) (*RegexpValidator, error) {
	re, err := regexp.Compile(s)
	if err != nil {
		return nil, err
	}
	return &RegexpValidator{Re: re}, nil
}

func (rv *RegexpValidator) acsepts(kind reflect.Kind) bool {
	return kind == reflect.String
}

func (rv *RegexpValidator) validateValue(item reflect.Value) error {
	ok := rv.Re.MatchString(item.String())
	if !ok {
		return NewErrMatch(rv.Re.String())
	}
	return nil
}
