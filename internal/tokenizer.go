package internal

import (
	"regexp"
	"strings"
)

func Tokenize(text string) []string {
	re := regexp.MustCompile(`\w+`)
	return re.FindAllString(strings.ToUpper(text), -1)
}
