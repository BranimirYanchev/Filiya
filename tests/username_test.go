package tests

import (
	"testing"

	"github.com/Marionvd/filia-project-backend/internal/username"
)

func TestNormalizeUsernameBase(t *testing.T) {
	cases := []struct {
		name     string
		fullName string
		expected string
	}{
		{
			name:     "Space Separated Names",
			fullName: "John Doe",
			expected: "john.doe",
		},
		{
			name:     "Multiple Separators",
			fullName: "  John   Doe-Smith  ",
			expected: "john.doe.smith",
		},
		{
			name:     "Cyrillic Name",
			fullName: "Иван Петров",
			expected: "иван.петров",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if actual := username.NormalizeBase(c.fullName); actual != c.expected {
				t.Fatalf("expected %q, got %q", c.expected, actual)
			}
		})
	}
}

func TestGenerateUniqueUsername(t *testing.T) {
	existing := map[string]bool{
		"john.doe":    true,
		"john.doe.10": true,
	}

	usernameValue, err := username.Generate("John Doe", func(candidate string) (bool, error) {
		return existing[candidate], nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if usernameValue != "john.doe.11" {
		t.Fatalf("expected john.doe.11, got %q", usernameValue)
	}
}
