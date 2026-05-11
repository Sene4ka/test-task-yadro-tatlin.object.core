package counter

import (
	"bufio"
	"io"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var russianTitle = cases.Title(language.Russian)

func CountNames(r io.Reader, preserveCase bool) (map[string]int, error) {
	result := make(map[string]int)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if name == "" {
			continue
		}
		if !preserveCase {
			name = russianTitle.String(strings.ToLower(name))
		}
		result[name]++
	}
	return result, scanner.Err()
}
