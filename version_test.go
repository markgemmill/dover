package main

import "fmt"
import "testing"
import "github.com/stretchr/testify/assert"

func assertVersion(t *testing.T, v *Version, major string, minor string, patch string, release string, build string) {
	msg := "Version.%s should be %s."
	assert.Equal(t, v.major, major, fmt.Sprintf(msg, "major", major))
	assert.Equal(t, v.minor, minor, fmt.Sprintf(msg, "minor", minor))
	assert.Equal(t, v.patch, patch, fmt.Sprintf(msg, "patch", patch))
	assert.Equal(t, v.release, release, fmt.Sprintf(msg, "release", release))
	assert.Equal(t, v.build, build, fmt.Sprintf(msg, "build", build))
}

func TestNewVersion(t *testing.T) {
	v := NewVersion([]string{"0", "1", "2", "", ""})
	assertVersion(t, v, "0", "1", "2", "", "0")
}

func TestVersionBumpMajor(t *testing.T) {
	v := NewVersion([]string{"0", "1", "2", "", ""})
	v2 := v.bumpMajor()
	assertVersion(t, &v2, "1", "0", "0", "", "0")
}

func TestVersionBumpMinor(t *testing.T) {
	v := NewVersion([]string{"0", "1", "2", "", ""})
	v2 := v.bumpMinor()
	assertVersion(t, &v2, "0", "2", "0", "", "0")
}

func TestVersionBumpPatch(t *testing.T) {
	v := NewVersion([]string{"0", "1", "2", "", ""})
	v2 := v.bumpPatch()
	assertVersion(t, &v2, "0", "1", "3", "", "0")
}

func TestVersionBumpReleaseToProd(t *testing.T) {
	v := NewVersion([]string{"0", "1", "2", "alpha", "1"})
	v2 := v.bumpReleaseToProd()
	assertVersion(t, &v2, "0", "1", "2", "", "0")
}

func TestVersionBumpRelease(t *testing.T) {
	var tests = []struct {
		release, expected string
	}{
		{"", "dev"},
		{"dev", "alpha"},
		{"alpha", "beta"},
		{"beta", "rc"},
		{"rc", ""},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("Test bumping from %s to %s", tt.release, tt.expected)
		t.Run(testname, func(t *testing.T) {
			v := NewVersion([]string{"0", "1", "2", tt.release, ""})
			v2 := v.bumpRelease()
			assertVersion(t, &v2, "0", "1", "2", tt.expected, "0")
		})
	}
}
