package pawnapi

import "testing"

func TestValidateVersion(t *testing.T) {
	valid := []string{"0.3.7", "1.0.0", "1", "4.0.5757-SF", "2000"}
	for _, v := range valid {
		if err := validateVersion(v); err != nil {
			t.Errorf("validateVersion(%q) = %v, want nil", v, err)
		}
	}

	invalid := []string{"", "latest", "v1.0.0", "1..0"}
	for _, v := range invalid {
		if err := validateVersion(v); err == nil {
			t.Errorf("validateVersion(%q) = nil, want error", v)
		}
	}
}

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"1.0.0", "1.0.0", 0},
		{"0.3.7", "1.0.0", -1},
		{"1.0.0", "0.3.7", 1},
		{"1.0.0", "1.0.0-rc1", 0}, // suffix ignored for the numeric comparison
		{"1.2", "1.2.0", 0},
		{"1.10.0", "1.9.0", 1},
	}
	for _, c := range cases {
		if got := compareVersions(c.a, c.b); got != c.want {
			t.Errorf("compareVersions(%q, %q) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}
