package main

import (
	"errors"
	"github.com/logrusorgru/aurora"
	"github.com/pelletier/go-toml"
	"os"
	"regexp"
	"strings"
)
import "io/fs"
import "fmt"

func check(e any) {
	if e != nil {
		panic(e)
	}
}

func findConfigFile() (string, error) {
	root := "./"
	fileSystem := os.DirFS(root)
	entries, err := fs.Glob(fileSystem, ".dover")
	check(err)

	if len(entries) == 1 {
		return entries[0], nil
	}

	return "", errors.New("Could not find dover config.")
}

func readConfigFile(configFile string) ([]string, error) {
	cfg, err := toml.LoadFile(configFile)
	check(err)

	files := cfg.GetArray("dover.versioned_files").([]string)

	return files, nil
}

func getVersionedFiles() []string {
	cfgFile, err := findConfigFile()
	check(err)

	cfg, err := readConfigFile(cfgFile)
	check(err)

	return cfg

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
	rx , _ := regexp.Compile(VERSION_REGEX)
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

func parseVersionString() {

}

func bumpVersion() {

}

func updateVersionStringInSourceFile() {

}

func doNext(part string) {
	files := getVersionedFiles()
	allMatches := make([]*VersionMatch, 0)
	for _, file := range files {
		content := readVersionSourceFile(file)
		for _, match := range searchForVersionString(file, content) {
			allMatches = append(allMatches, match)
		}
	}

	for _, match := range allMatches {
		//fmt.Printf("version: %s\n", match.version.format("000.A0"))
		//fmt.Printf("version: %s\n", match.version.format("000-A.0"))
		//fmt.Printf("version: %s\n", match.version.format("000-a0"))
		//fmt.Printf("version: %s\n", match.version.format("000+A0"))
		//fmt.Printf("version: %s\n", match.version.format("000+A.0"))
		fmt.Printf("%s:%0*d  %s\n", aurora.Yellow(match.file), 3, aurora.Blue(match.line), aurora.BrightWhite(match.version.toString()).Bold())
	}

}

func doBump(part string) {
	fmt.Println("buming version...")
	files := getVersionedFiles()
	for _, file := range files {
		fmt.Println(file)
	}
}

func runCommand(cmd string, part string) {
	switch cmd {
	case "next":
		doNext(part)
	case "bump":
		doBump(part)
	}
}

func main() {
	var args = len(os.Args)
	switch args {
	case 1:
		fmt.Print("Fetch the project version...")
	case 3:
		var cmd = os.Args[1]
		var part = os.Args[2]
		//fmt.Printf("Run %s %s\n", cmd, part)
		runCommand(cmd, part)
	default:
		fmt.Print("Invalid arguments....\n")
	}
}
