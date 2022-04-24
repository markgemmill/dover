package main

import "os"
import "regexp"
import "fmt"
import "strings"

type Formatter struct {
	versionFormat    string
	releaseSeparator string
	releaseFormat    string
	buildSeparator   string
	buildFormat      string
}

func (f *Formatter) format(v *Version) string {
	output := []string{}
	switch len(f.versionFormat) {
	case 1:
		output = append(output, v.major)
	case 2:
		output = append(output, v.major)
		output = append(output, ".")
		output = append(output, v.minor)
	case 3:
		output = append(output, v.major)
		output = append(output, ".")
		output = append(output, v.minor)
		output = append(output, ".")
		output = append(output, v.patch)
	}

	if v.release != "" {

		output = append(output, f.releaseSeparator)

		switch f.releaseFormat {
		case "a":
			output = append(output, SHORT[v.release])
		case "A":
			output = append(output, LONG[v.release])
		}

		if f.buildFormat == "0" {
			output = append(output, f.buildSeparator)
			output = append(output, v.build)
		}

	}

	return strings.Join(output, "")
}

var FORMAT_REGEX string = `^(000)([^a-zA-ZA\d])?([aA])?([^a-zA-Z\d])?(0)?$`

func NewVersionFormater(versionFormatString string) *Formatter {
	/// The format string consists of 5 parts:
	///  The numeric version format 000. Periods are assumed and there must be 3 zeros.
	///  The release separator - could be anything or nothing as long as it's not alphanumeric and it's a single character
	///	 The release name - either a or A to indicate abbreviated or long name.
	///  The build separator - could be anything or nothing as long as its not alphanumeric and it's a single character.
	///  The build number - this is either 0 or nothing.

	rx, err := regexp.Compile(FORMAT_REGEX)
	check(err)

	match := rx.FindStringSubmatch(versionFormatString)

	if match == nil {
		fmt.Printf("Invalid version format: %s", versionFormatString)
		os.Exit(1)
	}

	format := Formatter{
		versionFormat:    match[1],
		releaseSeparator: match[2],
		releaseFormat:    match[3],
		buildSeparator:   match[4],
		buildFormat:      match[5],
	}

	return &format
}
