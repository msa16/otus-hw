package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	ErrArgumentNotStructure = errors.New("argument is not a struct")
	ErrInvalidLenTag        = errors.New("invalid len tag value")
	ErrInvalidMinTag        = errors.New("invalid min tag value")
	ErrInvalidMaxTag        = errors.New("invalid max tag value")
	ErrInvalidInTag         = errors.New("invalid in tag value")
	ErrInvalidTag           = errors.New("invalid tag")
	ErrInvalidRegExp        = errors.New("invalid regexp")
	ErrWrongLen             = errors.New("wrong len")
	ErrRegExpNotMatch       = errors.New("regexp not matched")
	ErrNotInSet             = errors.New("not in set")
	ErrMinFailed            = errors.New("less than minimum")
	ErrMaxFailed            = errors.New("more than maximum")
	ErrUnknownValidator     = errors.New("unknown validator")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, e := range v {
		sb.WriteString("field: ")
		sb.WriteString(e.Field)
		sb.WriteString(", error: ")
		sb.WriteString(e.Err.Error())
	}
	return sb.String()
}

/*
Функция должна валидировать публичные поля входной структуры на основе структурного тэга `validate`.

Функция может возвращать
- или программную ошибку, произошедшую во время валидации;
- или `ValidationErrors` - ошибку, являющуюся слайсом структур, содержащих имя поля и ошибку его валидации.
*/

func Validate(v interface{}) error {
	errors := make(ValidationErrors, 0)
	programError := validateImpl(v, &errors)
	if programError != nil {
		return programError
	}
	if len(errors) == 0 {
		return nil
	}
	return errors
}

func validateImpl(v interface{}, errors *ValidationErrors) error {
	reflVal := reflect.ValueOf(v)
	// работаем только со структурами
	if reflVal.Kind() != reflect.Struct {
		return ErrArgumentNotStructure
	}
	numFields := reflVal.NumField()
	// проходимся по всем полям структуры
	for i := 0; i < numFields; i++ {
		structField := reflVal.Type().Field(i)
		// проверяем только публичные поля
		if !structField.IsExported() {
			continue
		}
		// Если у поля нет структурных тэгов или нет тэга validate, то функция игнорирует его.
		validate := structField.Tag.Get("validate")
		if validate == "" {
			continue
		}

		fieldVal := reflVal.Field(i)
		validateTags := strings.Split(validate, "|")
		var programError error
		switch {
		case fieldVal.CanInt():
			programError = validateInt(validateTags, structField.Name, int(fieldVal.Int()), errors)
		case fieldVal.Type().Kind() == reflect.String:
			programError = validateString(validateTags, structField.Name, fieldVal.String(), errors)
		case fieldVal.Type().Kind() == reflect.Slice:
			programError = validateSlice(validateTags, structField.Name, fieldVal, errors)
		case fieldVal.Kind() == reflect.Struct && validate == "nested":
			// поддержка валидации вложенных по композиции структур
			programError = validateImpl(fieldVal.Interface(), errors)
		}
		if programError != nil {
			return programError
		}
	}
	return nil
}

func validateSlice(validateTags []string, fieldName string, fieldVal reflect.Value, errors *ValidationErrors) error {
	if fieldVal.Type().Elem().Kind() == reflect.String {
		for i := 0; i < fieldVal.Len(); i++ {
			programError := validateString(validateTags, fieldName, fieldVal.Index(i).String(), errors)
			if programError != nil {
				return programError
			}
		}
	} else if fieldVal.Type().Elem().Kind() == reflect.Int {
		for i := 0; i < fieldVal.Len(); i++ {
			programError := validateInt(validateTags, fieldName, int(fieldVal.Index(i).Int()), errors)
			if programError != nil {
				return programError
			}
		}
	}
	return nil
}

func validateString(validateTags []string, fieldName string, fieldValue string, errors *ValidationErrors) error {
	for _, tag := range validateTags {
		t := strings.Split(tag, ":")
		if len(t) != 2 {
			return fmt.Errorf("%w: %s", ErrInvalidTag, tag)
		}
		switch t[0] {
		case "len":
			// длина строки должна быть ровно X символов
			expect, err := strconv.Atoi(t[1])
			if err != nil {
				return fmt.Errorf("%w: %s", ErrInvalidLenTag, t[1])
			}
			actual := utf8.RuneCountInString(fieldValue)
			if expect != actual {
				*errors = append(*errors, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("%w: expected %d, actual %d", ErrWrongLen, expect, actual),
				})
			}
		case "regexp":
			rg, err := regexp.Compile(t[1])
			if err != nil {
				return fmt.Errorf("%w: %s", ErrInvalidRegExp, t[1])
			}
			matched := rg.MatchString(fieldValue)
			if !matched {
				*errors = append(*errors, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("%w: %s", ErrRegExpNotMatch, t[1]),
				})
			}
		case "in":
			set := strings.Split(t[1], ",")
			found := false
			for _, s := range set {
				if s == fieldValue {
					found = true
					break
				}
			}
			if !found {
				*errors = append(*errors, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("%w: value %s, set %v", ErrNotInSet, fieldValue, set),
				})
			}
		default:
			return fmt.Errorf("%w: %s", ErrUnknownValidator, t[0])
		}
	}
	return nil
}

func validateInt(validateTags []string, fieldName string, fieldValue int, errors *ValidationErrors) error {
	for _, tag := range validateTags {
		t := strings.Split(tag, ":")
		if len(t) != 2 {
			return fmt.Errorf("%w: %s", ErrInvalidTag, tag)
		}
		switch t[0] {
		case "min":
			min, err := strconv.Atoi(t[1])
			if err != nil {
				return fmt.Errorf("%w: %s", ErrInvalidMinTag, t[1])
			}
			if fieldValue < min {
				*errors = append(*errors, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("%w: min %d, actual value %d", ErrMinFailed, min, fieldValue),
				})
			}
		case "max":
			max, err := strconv.Atoi(t[1])
			if err != nil {
				return fmt.Errorf("%w: %s", ErrInvalidMaxTag, t[1])
			}
			if fieldValue > max {
				*errors = append(*errors, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("%w: max %d, actual value %d", ErrMaxFailed, max, fieldValue),
				})
			}
		case "in":
			return validateIntIn(t[1], fieldName, fieldValue, errors)
		default:
			return fmt.Errorf("%w: %s", ErrUnknownValidator, t[0])
		}
	}
	return nil
}

func validateIntIn(inStr string, fieldName string, fieldValue int, errors *ValidationErrors) error {
	set := strings.Split(inStr, ",")
	found := false
	for _, s := range set {
		inValue, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidInTag, s)
		}
		if inValue == fieldValue {
			found = true
			break
		}
	}
	if !found {
		*errors = append(*errors, ValidationError{
			Field: fieldName,
			Err:   fmt.Errorf("%w: value %d, set %v", ErrNotInSet, fieldValue, set),
		})
	}
	return nil
}
