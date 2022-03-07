package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

type ValidationRule struct {
	name  string
	value string
}

var (
	ErrInvalidInputStruct      = errors.New("received interface is not a struct")
	ErrWrongValidationRule     = errors.New("wrong validation rule provided")
	ErrWrongNumericValue       = errors.New("numeric value is not recognized")
	ErrInvalidRegexpPattern    = errors.New("invalid regular expression pattern")
	ErrInvalidStringLength     = errors.New("invalid string length detected")
	ErrStringDoesntMatchRegexp = errors.New("string does not match the regexp pattern")
	ErrNumberLessThanRequired  = errors.New("number less than required limit")
	ErrNumberMoreThanRequired  = errors.New("number more than required limit")
	ErrValueNotInEnumList      = errors.New("number not in required enum list")
)

const (
	ruleLength = "len"
	ruleRegexp = "regexp"
	ruleMin    = "min"
	ruleMax    = "max"
	ruleIn     = "in"
)

func (v ValidationErrors) Error() string {
	errorList := make([]string, len(v))
	for i, valErr := range v {
		errorList[i] = fmt.Sprintf("Field `%s` contains an error: %s", valErr.Field, valErr.Err)
	}

	return strings.Join(errorList, ", ")
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return ErrInvalidInputStruct
	}

	var valErrors ValidationErrors
	valType := val.Type()

	for i := 0; i < valType.NumField(); i++ {
		field := valType.Field(i)
		fieldValue := val.Field(i)

		tag := field.Tag.Get("validate")
		if len(tag) == 0 {
			continue
		}

		validators := strings.Split(tag, "|")
		for _, validator := range validators {
			ruleStr := strings.Split(validator, ":")
			rule, err := getValidationRule(ruleStr)
			if err != nil {
				return err
			}

			err = checkValidationError(field, fieldValue, *rule)
			if err != nil {
				valErrors = append(valErrors, ValidationError{Field: field.Name, Err: err})
			}
		}
	}

	if len(valErrors) == 0 {
		return nil
	}

	return valErrors
}

func checkValidationError(field reflect.StructField, val reflect.Value, rule ValidationRule) error {
	var err error

	if field.Type.Kind() == reflect.Slice {
		items := reflect.ValueOf(val.Interface())
		for i := 0; i < items.Len(); i++ {
			item := items.Index(i)
			itemVal := reflect.Indirect(item)
			err = validateField(itemVal, rule.name, rule.value)
			if err != nil {
				break
			}
		}
	} else {
		err = validateField(val, rule.name, rule.value)
	}

	return err
}

func validateField(fieldValue reflect.Value, ruleName string, ruleValue string) error {
	var err error
	switch ruleName {
	case ruleLength:
		err = validateStringLength(fieldValue.String(), ruleValue)
	case ruleRegexp:
		err = validateMatchRegexp(fieldValue.String(), ruleValue)
	case ruleMin:
		err = validateMinValue(int(fieldValue.Int()), ruleValue)
	case ruleMax:
		err = validateMaxValue(int(fieldValue.Int()), ruleValue)
	case ruleIn:
		valueKind := fieldValue.Kind()
		if valueKind == reflect.String {
			err = validateEnumString(fieldValue.String(), ruleValue)
		} else if valueKind == reflect.Int {
			err = validateEnumInt(int(fieldValue.Int()), ruleValue)
		}
	}

	return err
}

func getValidationRule(rule []string) (*ValidationRule, error) {
	if ok := len(rule[0]) > 0; !ok {
		return nil, ErrWrongValidationRule
	}

	if ok := len(rule[1]) > 0; !ok {
		return nil, ErrWrongValidationRule
	}

	return &ValidationRule{
		name:  rule[0],
		value: rule[1],
	}, nil
}

func validateStringLength(str string, strLen string) error {
	intLen, err := strconv.Atoi(strLen)
	if err != nil {
		return ErrWrongNumericValue
	}

	if len(str) != intLen {
		return ErrInvalidStringLength
	}

	return nil
}

func validateMatchRegexp(str string, pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return ErrInvalidRegexpPattern
	}

	if re.Match([]byte(str)) {
		return nil
	}

	return ErrStringDoesntMatchRegexp
}

func validateMinValue(numVal int, minValueStr string) error {
	minValue, err := strconv.Atoi(minValueStr)
	if err != nil {
		return ErrWrongNumericValue
	}

	if numVal < minValue {
		return ErrNumberLessThanRequired
	}

	return nil
}

func validateMaxValue(numVal int, maxValueStr string) error {
	maxValue, err := strconv.Atoi(maxValueStr)
	if err != nil {
		return ErrWrongNumericValue
	}

	if numVal > maxValue {
		return ErrNumberMoreThanRequired
	}

	return nil
}

func validateEnumString(fieldVal string, enumStr string) error {
	enum := strings.Split(enumStr, ",")
	if len(enum) == 0 {
		return nil
	}

	for _, v := range enum {
		if fieldVal == v {
			return nil
		}
	}

	return ErrValueNotInEnumList
}

func validateEnumInt(fieldVal int, enumStr string) error {
	enum := strings.Split(enumStr, ",")
	if len(enum) == 0 {
		return nil
	}

	for _, v := range enum {
		elem, err := strconv.Atoi(v)
		if err != nil {
			return ErrWrongNumericValue
		}

		if fieldVal == elem {
			return nil
		}
	}

	return ErrValueNotInEnumList
}
