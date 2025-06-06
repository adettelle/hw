package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	escapeMode := false
	r := []rune(s)

	var sb strings.Builder

	for i := 0; i < len(r); i++ {
		if !escapeMode && unicode.IsDigit(r[i]) { // первый элемент цифра
			return "", ErrInvalidString
		}

		if i == len(r)-1 {
			if string(r[i]) == `\` { // последний элемент `\`
				return "", ErrInvalidString
			}
			// последний элемент любой другой
			sb.WriteString(string(r[i]))
			return sb.String(), nil
		}
		// ----------------- не последний элемент -----------------

		if !escapeMode && string(r[i]) == `\` {
			escapeMode = true
			continue
		}

		if escapeMode && (!unicode.IsDigit(r[i]) && string(r[i]) != `\`) {
			return "", ErrInvalidString
		}

		if unicode.IsDigit(r[i+1]) { // `?4` следующий элемент - цифра
			n, err := strconv.Atoi(string(r[i+1]))
			if err != nil {
				return "", err
			}

			sb.WriteString(strings.Repeat(string(r[i]), n))

			i++
		} else { // `ab` следующий элемент - не цифра
			sb.WriteString(string(r[i]))
			escapeMode = false
		}
	}
	return sb.String(), nil
}
