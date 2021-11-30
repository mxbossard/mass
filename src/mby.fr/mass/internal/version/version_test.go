package version

import (
	"errors"
	"testing"

	"github.com/Masterminds/semver/v3"
)

func TestNextPatch(t *testing.T) {
	cases := []struct {
		in, want string; err error
	}{
		{"0.1.2", "0.1.3", nil},
		{"1.0.7-dev", "1.0.7", nil},
		{"1.0.3-rc1", "1.0.3", nil},
		{"1", "1.0.1", nil},
		{"1.1", "1.1.1", nil},
		{"", "", semver.ErrInvalidSemVer},
		{"foo", "", semver.ErrInvalidSemVer},
	}
	for _, c := range cases {
		got, err := NextPatch(c.in)
		if got != c.want {
			t.Errorf("NextPatch(%q) == %q, want %q", c.in, got, c.want)
		}
		if !errors.Is(err, c.err) {
			t.Errorf("NextPatch(%q) error == %q, expected %q", c.in, err, c.err)
		}
	}
}
