package main

import (
	"strconv"
)

const VERSION_REGEX = `version[^ :=]* ?[:=] ? ["']?(?P<major>\d+)(\.(?P<minor>\d+))(\.(?P<patch>\d+))?([\.-](?P<release>[a-z]+)([\.-]?(?P<build>\d+))?)?["']?`
//var VERSION_REGEX *regexp.Regexp

var RELEASE = [4]string{"dev", "alpha", "beta", "rc"}
var RELEASES = map[string]string{ "dev": "d", "alpha": "a", "beta": "b", "rc": "rc", "d": "dev", "a": "alpha", "b": "beta"}
var LONG = map[string]string{ "dev": "dev", "alpha": "alpha", "beta": "beta", "rc": "rc", "d": "dev", "a": "alpha", "b": "beta"}
var SHORT = map[string]string{ "dev": "d", "alpha": "a", "beta": "b", "rc": "rc", "d": "d", "a": "a", "b": "b"}


//func Init() {
//	VERSION_REGEX, _ = regexp.Compile(VERSION)
//}


func parseRegexResults(match []string) []string {
	return []string{
		match[1],
		match[3],
		match[5],
		match[7],
		match[9],
	}
}


func nextRelease(currentRelease string) string {
	 index := -1
	 for i, release := range RELEASE {
	 	if release == currentRelease {
	 		index = i
		}
	 }
	 index += 1
	 if index + 1 > len(RELEASE) {
	 	return ""
	 }
	 return RELEASE[index]
}


type Version struct {
	major   string
	minor   string
	patch   string
	release string
	build   string
}

func (v *Version) format(fmtString string) string {
	f:= NewVersionFormater(fmtString)
	return f.format(v)
}

func (v *Version) toString() string {
	return v.format("000-A.0")
}

func (v *Version) bumpMajor() *Version {
	major, err := strconv.Atoi(v.major)
	check(err)

	major = major + 1
	nv := Version{
		major: strconv.Itoa(major),
		minor: "0",
		patch: "0",
		release: "",
		build: "0",
	}
	return &nv

}

func (v *Version) bumpMinor() *Version {
	minor, err := strconv.Atoi(v.minor)
	check(err)

	minor = minor + 1
	nv := Version{
		major: v.major,
		minor: strconv.Itoa(minor),
		patch: "0",
		release: "",
		build: "0",
	}
	return &nv
}

func (v *Version) bumpPatch() *Version {
	patch, err := strconv.Atoi(v.patch)
	check(err)

	patch = patch + 1
	nv := Version{
		major: v.major,
		minor: v.minor,
		patch: strconv.Itoa(patch),
		release: "",
		build: "0",
	}
	return &nv
}

func (v *Version) bumpRelease() *Version {
	release := nextRelease(v.release)

	nv := Version{
		major: v.major,
		minor: v.minor,
		patch: v.patch,
		release: release,
		build: "0",
	}
	return &nv

}

func (v *Version) bumpReleaseToProd() *Version {
	nv := Version{
		major: v.major,
		minor: v.minor,
		patch: v.patch,
		release: "",
		build: "0",
	}
	return &nv
}


func (v *Version) bumpBuild() *Version {

	build, err := strconv.Atoi(v.build)
	check(err)

	build = build + 1
	nv := Version{
		major: v.major,
		minor: v.minor,
		patch: v.patch,
		release: v.release,
		build: strconv.Itoa(build),
	}
	return &nv
}

func (v *Version) bump(part string) *Version {
	var newVers *Version
	switch part {
	case "major":
		newVers = v.bumpMajor()
	case "minor":
		newVers = v.bumpMinor()
	case "patch":
		newVers = v.bumpPatch()
	case "release":
		newVers = v.bumpRelease()
	case "prod":
		newVers = v.bumpReleaseToProd()
	case "build":
		newVers = v.bumpBuild()
	default:
		newVers = v
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
