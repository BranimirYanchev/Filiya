package username

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

func NormalizeBase(fullname string) string {
	var builder strings.Builder
	lastWasSeparator := false

	for _, r := range strings.ToLower(strings.TrimSpace(fullname)) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			builder.WriteRune(r)
			lastWasSeparator = false
		case unicode.IsSpace(r) || r == '-' || r == '_' || r == '.':
			if builder.Len() > 0 && !lastWasSeparator {
				builder.WriteRune('.')
				lastWasSeparator = true
			}
		}
	}

	username := strings.Trim(builder.String(), ".")
	if username == "" {
		return "user"
	}

	return username
}

func Generate(fullname string, exists func(string) (bool, error)) (string, error) {
	base := NormalizeBase(fullname)

	used, err := exists(base)
	if err != nil {
		return "", err
	}
	if !used {
		return base, nil
	}

	for suffix := 10; suffix <= 999; suffix++ {
		candidate := fmt.Sprintf("%s.%d", base, suffix)
		used, err = exists(candidate)
		if err != nil {
			return "", err
		}
		if !used {
			return candidate, nil
		}
	}

	return "", errors.New("could not generate unique username")
}
