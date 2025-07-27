package app

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// TODO: combines these strings
const (
	JUST_VERSION  = `(?P<major>\d+)(\.(?P<minor>\d+))(\.(?P<patch>\d+))?([\.\-\+](?P<release>[a-z]+)([\.-]?(?P<build>\d+))?)?`
	VERSION_REGEX = `(version|VERSION|Version)[^ :=]* ?[:=]? ? ["']?(?P<major>\d+)(\.(?P<minor>\d+))(\.(?P<patch>\d+))?([\.\-\+](?P<release>[a-z]+)([\.-]?(?P<build>\d+))?)?["']?`
)

var (
	RELEASE  = []string{"dev", "alpha", "beta", "rc"}
	RELEASES = map[string]string{"dev": "d", "alpha": "a", "beta": "b", "rc": "rc", "d": "dev", "a": "alpha", "b": "beta"}
	SHORT    = map[string]string{"dev": "d", "alpha": "a", "beta": "b", "rc": "rc", "d": "d", "a": "a", "b": "b"}
	LONG     = map[string]string{"dev": "dev", "alpha": "alpha", "beta": "beta", "rc": "rc", "d": "dev", "a": "alpha", "b": "beta"}
	PARTS    = [6]string{"major", "minor", "patch", "release", "prod", "build"}
)

type VersionFinder struct {
	rx regexp.Regexp
}

func (vf *VersionFinder) parseRegexResults(match []string) []string {
	return []string{
		match[2],  // major
		match[4],  // minor
		match[6],  // patch
		match[8],  // release
		match[10], // build
	}
}

func (vf *VersionFinder) Find(line string) (Version, bool) {
	match := vf.rx.FindStringSubmatch(line)
	if match != nil {
		return *NewVersion(vf.parseRegexResults(match)), true
	}
	return Version{
		major:   "",
		minor:   "",
		patch:   "",
		release: "",
		build:   "",
	}, false
}

func NewVersionFinder() *VersionFinder {
	_rx, _ := regexp.Compile(VERSION_REGEX)
	vf := VersionFinder{
		rx: *_rx,
	}
	return &vf
}

func nextRelease(currentRelease string) string {
	index := IndexOf(&RELEASE, currentRelease)
	index += 1
	if index+1 > len(RELEASE) {
		return ""
	}
	return RELEASE[index]
}

func validateReleaseOrder(currentRelease string, requestedRelease string) error {
	currentIndex := IndexOf(&RELEASE, currentRelease)
	requestedIndex := IndexOf(&RELEASE, requestedRelease)
	if requestedIndex < currentIndex {
		msg := fmt.Sprintf("Invalid release order requested. `%s` comes before the current release `%s`.", requestedRelease, currentRelease)
		return errors.New(msg)
	}
	return nil
}

type Version struct {
	major   string
	minor   string
	patch   string
	release string
	build   string
}

func (v *Version) copy() Version {
	nv := Version{
		major:   v.major,
		minor:   v.minor,
		patch:   v.patch,
		release: v.release,
		build:   v.build,
	}
	return nv
}

func (v *Version) format(fmtString string) string {
	f := NewVersionFormater(fmtString)
	return f.format(v)
}

func (v *Version) toString() string {
	return v.format("000-A.0")
}

func (v *Version) bumpMajor() Version {
	major, err := strconv.Atoi(v.major)
	check(err)

	major = major + 1
	nv := Version{
		major:   strconv.Itoa(major),
		minor:   "0",
		patch:   "0",
		release: "",
		build:   "0",
	}
	return nv
}

func (v *Version) bumpMinor() Version {
	minor, err := strconv.Atoi(v.minor)
	check(err)

	minor = minor + 1
	nv := Version{
		major:   v.major,
		minor:   strconv.Itoa(minor),
		patch:   "0",
		release: "",
		build:   "0",
	}
	return nv
}

func (v *Version) bumpPatch() Version {
	patch, err := strconv.Atoi(v.patch)
	check(err)

	patch = patch + 1
	nv := Version{
		major:   v.major,
		minor:   v.minor,
		patch:   strconv.Itoa(patch),
		release: "",
		build:   "0",
	}
	return nv
}

func (v *Version) setPreRelease(release string) Version {
	err := validateReleaseOrder(v.release, release)
	ExitOnError(err)

	nv := v.copy()
	if nv.release != release {
		nv.release = release
		nv.build = "0"
	}
	return nv
}

func (v *Version) bumpRelease() Version {
	release := nextRelease(v.release)

	nv := Version{
		major:   v.major,
		minor:   v.minor,
		patch:   v.patch,
		release: release,
		build:   "0",
	}
	return nv
}

func (v *Version) bumpReleaseToProd() Version {
	nv := Version{
		major:   v.major,
		minor:   v.minor,
		patch:   v.patch,
		release: "",
		build:   "0",
	}
	return nv
}

func (v *Version) bumpBuild() Version {
	build, err := strconv.Atoi(v.build)
	check(err)

	build = build + 1
	nv := Version{
		major:   v.major,
		minor:   v.minor,
		patch:   v.patch,
		release: v.release,
		build:   strconv.Itoa(build),
	}
	return nv
}

func (v *Version) bump(part string, preRelease string) Version {
	newVers := v.copy()

	switch part {
	case "major":
		newVers = newVers.bumpMajor()
	case "minor":
		newVers = newVers.bumpMinor()
	case "patch":
		newVers = newVers.bumpPatch()
	}

	switch preRelease {
	case "pre-release":
		newVers = newVers.bumpRelease()
	case "dev":
		newVers = newVers.setPreRelease(preRelease)
	case "alpha":
		newVers = newVers.setPreRelease(preRelease)
	case "beta":
		newVers = newVers.setPreRelease(preRelease)
	case "rc":
		newVers = newVers.setPreRelease(preRelease)
	case "release":
		newVers = newVers.bumpReleaseToProd()
	}

	if newVers.release != "" && (part == "build" || v.release == preRelease) {
		newVers = newVers.bumpBuild()
	}

	return newVers
}

func (v *Version) equals(other *Version) bool {
	return v.toString() == other.toString()
}

func defaultZeroStr(input string) string {
	if input == "" {
		return "0"
	}
	return input
}

func NewVersion(match []string) *Version {
	v := Version{
		major:   defaultZeroStr(match[0]),
		minor:   defaultZeroStr(match[1]),
		patch:   defaultZeroStr(match[2]),
		release: match[3],
		build:   defaultZeroStr(match[4]),
	}
	return &v
}
