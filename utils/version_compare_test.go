package utils

import "testing"

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		a, b     string
		expected int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.0", "1.1.0", -1},
		{"1.2.0", "1.1.9", 1},
		{"1.0", "1.0.1", -1},
	}
	for _, c := range cases {
		got := CompareVersions(c.a, c.b)
		if got != c.expected {
			t.Errorf("CompareVersions(%s, %s) = %d, expected %d", c.a, c.b, got, c.expected)
		}
	}
}
