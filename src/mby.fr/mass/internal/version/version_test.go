package version

import (
	"testing"
)

func TestBumpPatch(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"0.0.1", "0.0.2"},
		{"1.0.0-rc1", "1.0.1-rc1"},
		{"", ""},
	}
	for _, c := range cases {
		got := BumpPatch(c.in)
		if got != c.want {
			t.Errorf("BumpPatch(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
