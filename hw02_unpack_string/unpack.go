package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrInvalidString  = errors.New("invalid string")
	ErrStringNotAlnum = errors.New("alphanumeric only allowed")
)

func Unpack(str string) (string, error) {
	if len(str) == 0 {
		return "", nil
	}

	_, validationError := validate(str)
	if nil != validationError {
		return str, validationError
	}

	prev := str[0]
	var b strings.Builder

	for i := 1; i < len(str); i++ {
		current := str[i]
		if unicode.IsDigit(rune(current)) {
			count, _ := strconv.Atoi(string(current))
			repeated := strings.Repeat(string(prev), count)
			b.WriteString(repeated)
			prev = str[i+1]
			i++
		} else {
			b.WriteString(string(prev))
			prev = current
		}
	}

	last := string(str[len(str)-1])
	b.WriteString(last)

	return b.String(), nil
}

func validate(str string) (string, error) {
	prev := str[0]
	if unicode.IsDigit(rune(prev)) {
		return str, ErrInvalidString
	}

	for i := 0; i < len(str); i++ {
		if i == 0 {
			continue
		}

		current := str[i]

		alphanumeric := unicode.IsDigit(rune(current)) || unicode.IsLetter(rune(current))
		if !alphanumeric {
			return str, ErrStringNotAlnum
		}

		if unicode.IsDigit(rune(current)) && unicode.IsDigit(rune(prev)) {
			return str, ErrInvalidString
		}

		prev = current
	}

	return str, nil
}
