package main

import "os"
import "regexp"
import "fmt"
import "strings"

/// Format
///  000.a0
///  000+a.0
///  000-A0

type Formater struct {
	versionFormat string
	releaseSeparator string
	releaseFormat string
	buildSeparator string
	buildFormat string
}

func (f *Formater) format(v *Version) string {
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
	//fmt.Println("format output: ", output)

	return strings.Join(output, "")
}


var FORMAT_REGEX string = `^(000)([^a-zA-ZA\d])?([aA])?([^a-zA-Z\d])?(0)?$`


func NewVersionFormater(versionFormatString string) *Formater{
	rx, err := regexp.Compile(FORMAT_REGEX)
	check(err)

	match := rx.FindStringSubmatch(versionFormatString)
	//fmt.Println("format version match: ", match)

	if match == nil {
		fmt.Printf("Invalid version format: %s", versionFormatString)
		os.Exit(1)
	}

	format := Formater{
		versionFormat: match[1],
		releaseSeparator: match[2],
		releaseFormat: match[3],
		buildSeparator: match[4],
		buildFormat: match[5],
	}
	//fmt.Printf("format formater: %v\n", format)
	return &format
}
