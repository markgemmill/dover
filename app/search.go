package app

import (
	"errors"
	"os"
	"strconv"
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

func searchForVersionString(file string, lines []int, fileContent []string) []*VersionMatch {
	lineMatches := make([]*VersionMatch, 0)
	finder := NewVersionFinder()
	for index, line := range fileContent {
		if len(lines) > 0 && IndexOf(&lines, index) == -1 {
			continue
		}
		v, found := finder.Find(line)
		if found {
			vm := newVersionMatch(file, index, &v)
			lineMatches = append(lineMatches, vm)
		}
	}
	return lineMatches
}

func parseVersionedFileConfig(filePath string) (string, []int) {
	lines := make([]int, 0)
	filePath, lineNotation := splitFileAndLineNotation(filePath)
	if lineNotation == "" {
		return filePath, lines
	}
	for _, value := range strings.Split(lineNotation, ",") {
		lineInt, err := strconv.Atoi(value)
		if err != nil {
			ExitOnError(errors.New("invalid versioned file:line notation in configuration file"))
		}
		lines = append(lines, lineInt)
	}
	return filePath, lines
}

func getAllVersionStringMatches(files []string) *[]*VersionMatch {
	allMatches := make([]*VersionMatch, 0)
	for _, file := range files {
		filePath, lines := parseVersionedFileConfig(file)
		content := readVersionSourceFile(filePath)
		for _, match := range searchForVersionString(filePath, lines, content) {
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
