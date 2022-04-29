package app

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func configValues() (ConfigValues, error) {

	possibleConfigFiles := map[string]configParser{
		".dover":         getTomlConfigValues,
		"pyproject.toml": getTomlConfigValues,
		"package.json":   getJSONConfigValues,
	}

	var cfg ConfigValues

	for fileName, configParser := range possibleConfigFiles {
		cfgFile, err := findConfigFile(fileName)
		if err != nil {
			continue
		}

		cfg, err := configParser(cfgFile)
		if err != nil {
			fmt.Printf("%s: %s", fileName, err)
			continue
		}
		return cfg, nil
	}

	return cfg, errors.New("Unable to find dover configuration.")

}

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
	rx, _ := regexp.Compile(VERSION_REGEX)
	for index, line := range fileContent {
		match := rx.FindStringSubmatch(line)
		if match != nil {
			v := NewVersion(parseRegexResults(match))
			vm := newVersionMatch(file, index, v)
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
