package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var sb strings.Builder
	var prev rune
	escapeMode := false
	for _, r := range str {
		digit, err := strconv.Atoi(string(r))
		if escapeMode {
			// были в состоянии экранирования
			escapeMode = false
			if err != nil && r != '\\' {
				// экранировать можно только цифру или обратный слеш
				return "", ErrInvalidString
			}
			prev = r
			continue
		}

		if err == nil {
			// очередной символ - цифра
			if prev == 0 {
				// перед цифрой не было символа
				return "", ErrInvalidString
			}
			sb.WriteString(strings.Repeat(string(prev), digit))
			prev = 0
		} else {
			if prev > 0 {
				sb.WriteRune(prev)
			}
			escapeMode = r == '\\'
			prev = r
		}
	}
	if escapeMode {
		// экранирование не было закончено
		return "", ErrInvalidString
	}
	if prev > 0 {
		sb.WriteRune(prev)
	}
	return sb.String(), nil
}
