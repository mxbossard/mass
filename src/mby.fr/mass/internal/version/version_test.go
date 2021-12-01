package version

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNextPatch(t *testing.T) {
	cases := []struct {
		in, want, err string
	}{
		{"0.1.2", "0.1.3", ""},
		{"1.0.7-dev", "1.0.7", ""},
		{"1.0.3-rc1", "1.0.3", ""},
		{"1", "1.0.1", ""},
		{"1.1", "1.1.1", ""},
		{"", "", "Invalid Semantic Version"},
		{"foo", "", "Invalid Semantic Version"},
	}
	for i, c := range cases {
		got, err := NextPatch(c.in)
		assert.Equal(t, c.want, got, "case #%d should be equal", i)
		if c.err != "" {
			assert.EqualError(t, err, c.err)
		}
	}
}

func TestNextMinor(t *testing.T) {
	cases := []struct {
		in, want, err string
	}{
		{"0.1.2", "0.2.0", ""},
		{"1.0.7-dev", "1.1.0", ""},
		{"1.0.3-rc1", "1.1.0", ""},
		{"1", "1.1.0", ""},
		{"1.1", "1.2.0", ""},
		{"", "", "Invalid Semantic Version"},
		{"foo", "", "Invalid Semantic Version"},
	}
	for i, c := range cases {
		got, err := NextMinor(c.in)
		assert.Equal(t, c.want, got, "case #%d should be equal", i)
		if c.err != "" {
			assert.EqualError(t, err, c.err)
		}
	}
}

func TestNextMajor(t *testing.T) {
	cases := []struct {
		in, want, err string
	}{
		{"0.1.2", "1.0.0", ""},
		{"1.0.7-dev", "2.0.0", ""},
		{"1.0.3-rc1", "2.0.0", ""},
		{"1", "2.0.0", ""},
		{"1.1", "2.0.0", ""},
		{"", "", "Invalid Semantic Version"},
		{"foo", "", "Invalid Semantic Version"},
	}
	for i, c := range cases {
		got, err := NextMajor(c.in)
		assert.Equal(t, c.want, got, "case #%d should be equal", i)
		if c.err != "" {
			assert.EqualError(t, err, c.err)
		}
	}
}

func TestIsDev(t *testing.T) {
	cases := []struct {
		in string
		want bool
		err string
	}{
		{"0.1.2", false, ""},
		{"1.0.7-dev", true, ""},
		{"1.0.3-rc1", false, ""},
		{"1", false, ""},
		{"1.1", false, ""},
		{"", false, "Invalid Semantic Version"},
		{"foo", false, "Invalid Semantic Version"},
	}
	for i, c := range cases {
		got, err := IsDev(c.in)
		assert.Equal(t, c.want, got, "case #%d should be equal", i)
		if c.err != "" {
			assert.EqualError(t, err, c.err)
		}
	}
}

func TestIsRc(t *testing.T) {
	cases := []struct {
		in string
		want bool
		err string
	}{
		{"0.1.2", false, ""},
		{"1.0.7-dev", false, ""},
		{"1.0.3-rc1", true, ""},
		{"1", false, ""},
		{"1.1", false, ""},
		{"", false, "Invalid Semantic Version"},
		{"foo", false, "Invalid Semantic Version"},
	}
	for i, c := range cases {
		got, err := IsRc(c.in)
		assert.Equal(t, c.want, got, "case #%d should be equal", i)
		if c.err != "" {
			assert.EqualError(t, err, c.err)
		}
	}
}

func TestNextDev(t *testing.T) {
	cases := []struct {
		in string
		want string
		err string
	}{
		{"0.1.2", "0.1.3-dev", ""},
		{"1.0.7-dev", "1.0.7-dev", ""},
		{"1.0.3-rc1", "1.0.4-dev", ""},
		{"1", "1.0.1-dev", ""},
		{"1.1", "1.1.1-dev", ""},
		{"", "", "Invalid Semantic Version"},
		{"foo", "", "Invalid Semantic Version"},
	}
	for i, c := range cases {
		got, err := NextDev(c.in)
		assert.Equal(t, c.want, got, "case #%d should be equal", i)
		if c.err != "" {
			assert.EqualError(t, err, c.err)
		}
	}
}

func TestNextRc(t *testing.T) {
	cases := []struct {
		in string
		want string
		err string
	}{
		{"0.1.2", "0.1.3-rc1", ""},
		{"1.0.7-dev", "1.0.7-rc1", ""},
		{"1.0.3-rc1", "1.0.3-rc2", ""},
		{"1", "1.0.1-rc1", ""},
		{"1.1", "1.1.1-rc1", ""},
		{"", "", "Invalid Semantic Version"},
		{"foo", "", "Invalid Semantic Version"},
	}
	for i, c := range cases {
		got, err := NextRc(c.in)
		assert.Equal(t, c.want, got, "case #%d should be equal", i)
		if c.err != "" {
			assert.EqualError(t, err, c.err)
		}
	}
}
