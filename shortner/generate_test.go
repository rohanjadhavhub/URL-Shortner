package shortner

import (
	"strings"
	"testing"
)

func TestGenerateCode_Length(t *testing.T) {
	lengths := []int{0, 1, 7, 20}
	for _, l := range lengths {
		got := GenerateCode(l)
		if len(got) != l {
			t.Errorf("GenerateCode(%d) length = %d, want %d", l, len(got), l)
		}
	}
}

func TestGenerateCode_Charset(t *testing.T) {
	code := GenerateCode(100)
	for _, c := range code {
		if !strings.ContainsRune(charset, c) {
			t.Errorf("GenerateCode produced invalid character %q not in charset", c)
		}
	}
}

func TestGenerateCode_Randomness(t *testing.T) {
	const n = 1000
	seen := make(map[string]bool, n)
	for i := 0; i < n; i++ {
		code := GenerateCode(7)
		if seen[code] {
			t.Fatalf("GenerateCode produced a duplicate after %d calls: %q", i, code)
		}
		seen[code] = true
	}
}
