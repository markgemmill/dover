package app

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

func writeVersionUpdate(filePath string, lineNo int, newVersion string) {
	content, err := os.ReadFile(filePath)
	check(err)

	lines := strings.Split(string(content), "\n")
	lineTotal := len(lines)

	rx := regexp.MustCompile(JUST_VERSION)

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	ExitOnError(err)

	err = file.Truncate(0)
	ExitOnError(err)

	_, err = file.Seek(0, 0)
	ExitOnError(err)

	defer file.Close()

	writer := bufio.NewWriter(file)
	for index, line := range lines {
		lineSeparator := "\n"
		if index+1 == lineTotal {
			lineSeparator = ""
		}
		if index == lineNo {
			newLine := rx.ReplaceAllString(line, newVersion)
			_, err = writer.WriteString(newLine + lineSeparator)
			ExitOnError(err)
		} else {
			_, err = writer.WriteString(line + lineSeparator)
			ExitOnError(err)
		}
	}
	writer.Flush()

}
