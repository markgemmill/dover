package app

import (
	"os"
	"strings"
)

type VersionMatch struct {
	file    string
	line    int
	version *Version
}

func newVersionMatch(file string, line int, version *Version) *VersionMatch {
	vm := VersionMatch{
		file:    file,
		line:    line,
		version: version,
	}
	return &vm
}

func readVersionSourceFile(filePath string) []string {
	content, err := os.ReadFile(filePath)
	check(err)
	return strings.Split(string(content), "\n")
}

func searchForVersionString(file string, fileContent []string) []*VersionMatch {
	lineMatches := make([]*VersionMatch, 0)
	finder := NewVersionFinder()
	for index, line := range fileContent {
		v, found := finder.Find(line)
		if found {
			vm := newVersionMatch(file, index, &v)
			lineMatches = append(lineMatches, vm)
		}
	}
	return lineMatches
}

func getAllVersionStringMatches(files []string) *[]*VersionMatch {

	allMatches := make([]*VersionMatch, 0)
	for _, file := range files {
		content := readVersionSourceFile(file)
		for _, match := range searchForVersionString(file, content) {
			allMatches = append(allMatches, match)
		}
	}

	return &allMatches
}

func assertVersionMatchConsistency(matches *[]*VersionMatch) bool {
	var rootVersion *Version = (*matches)[0].version
	for _, m := range *matches {
		if !m.version.equals(rootVersion) {
			return false
		}
	}
	return true
}
