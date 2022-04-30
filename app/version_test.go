package app

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

func TestVersionCopy(t *testing.T) {
	v1 := NewVersion([]string{"0", "1", "2", "", ""})
	v2 := v1.copy()

	assert.Equal(t, v2.toString(), v1.toString())
	assert.True(t, v1.equals(&v2))
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

func TestVersionBumpBuild(t *testing.T) {
	v := NewVersion([]string{"0", "1", "2", "dev", "0"})
	v2 := v.bumpBuild()
	assertVersion(t, &v2, "0", "1", "2", "dev", "1")
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

func TestValidateReleaseOrder(t *testing.T) {
	// dev
	err := validateReleaseOrder("", "dev")
	assert.Nil(t, err)

	err = validateReleaseOrder("dev", "dev")
	assert.Nil(t, err)

	err = validateReleaseOrder("alpha", "dev")
	assert.NotNil(t, err)

	err = validateReleaseOrder("beta", "dev")
	assert.NotNil(t, err)

	err = validateReleaseOrder("rc", "dev")
	assert.NotNil(t, err)

	// alpha
	err = validateReleaseOrder("", "alpha")
	assert.Nil(t, err)

	err = validateReleaseOrder("dev", "alpha")
	assert.Nil(t, err)

	err = validateReleaseOrder("alpha", "alpha")
	assert.Nil(t, err)

	err = validateReleaseOrder("beta", "alpha")
	assert.NotNil(t, err)

	err = validateReleaseOrder("rc", "alpha")
	assert.NotNil(t, err)

	// beta
	err = validateReleaseOrder("", "beta")
	assert.Nil(t, err)

	err = validateReleaseOrder("dev", "beta")
	assert.Nil(t, err)

	err = validateReleaseOrder("alpha", "beta")
	assert.Nil(t, err)

	err = validateReleaseOrder("beta", "beta")
	assert.Nil(t, err)

	err = validateReleaseOrder("rc", "beta")
	assert.NotNil(t, err)

	// rc
	err = validateReleaseOrder("", "rc")
	assert.Nil(t, err)

	err = validateReleaseOrder("dev", "rc")
	assert.Nil(t, err)

	err = validateReleaseOrder("alpha", "rc")
	assert.Nil(t, err)

	err = validateReleaseOrder("beta", "rc")
	assert.Nil(t, err)

	err = validateReleaseOrder("rc", "rc")
	assert.Nil(t, err)

}

func assertNewVersion(t *testing.T, version *Version, part string, preRelease string, equals string) {
	v2 := version.bump(part, preRelease)
	assert.Equal(t, equals, v2.toString())

}

func TestVersionBump(t *testing.T) {
	v1 := NewVersion([]string{"0", "1", "2", "dev", "0"})

	assertNewVersion(t, v1, "major", "", "1.0.0")
	assertNewVersion(t, v1, "minor", "", "0.2.0")
	assertNewVersion(t, v1, "patch", "", "0.1.3")
	assertNewVersion(t, v1, "build", "", "0.1.2-dev.1")
	assertNewVersion(t, v1, "", "dev", "0.1.2-dev.1")
	assertNewVersion(t, v1, "", "alpha", "0.1.2-alpha.0")
	assertNewVersion(t, v1, "", "beta", "0.1.2-beta.0")
	assertNewVersion(t, v1, "", "rc", "0.1.2-rc.0")

	assertNewVersion(t, v1, "build", "dev", "0.1.2-dev.1")
	assertNewVersion(t, v1, "build", "alpha", "0.1.2-alpha.1")
	assertNewVersion(t, v1, "build", "beta", "0.1.2-beta.1")
	assertNewVersion(t, v1, "build", "rc", "0.1.2-rc.1")
	assertNewVersion(t, v1, "", "release", "0.1.2")

}
