package str

import (
	"unicode/utf8"
)

func CutStr(str string, n int) string {

	length := utf8.RuneCountInString(str)

	if length < n {
		return str
	} else {
		var substring []rune
		for i, char := range str {
			if i < n {
				substring = append(substring, char)
			}
		}
		return string(substring) + "..."
	}
}
