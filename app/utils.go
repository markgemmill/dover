package app

import (
	"fmt"
	"os"
	"strings"

	c "github.com/fatih/color"
)

func check(e any) {
	if e != nil {
		panic(e)
	}
}

func ExitOnError(e error) {
	if e != nil {
		fmt.Printf("%s\n", c.RedString(fmt.Sprintf("%s", e)))
		os.Exit(1)
	}
}

func setMax(currentValue int, maxValue *int) {
	if currentValue > *maxValue {
		*maxValue = currentValue
	}
}

func maxLength(items []string) int {
	maxLen := 0
	for _, item := range items {
		itemLen := len(item)
		if itemLen > maxLen {
			maxLen = itemLen
		}
	}
	return maxLen
}

func digitCount(number int) int {
	if number < 10 {
		return 1
	}
	if number < 100 {
		return 2
	}
	if number < 1000 {
		return 3
	}
	if number < 10000 {
		return 4
	}
	return 5
}

func getMaxColumnWidths(matches *[]*VersionMatch, format string) (int, int, int) {
	fileWidth := 0
	lineWidth := 0
	versionWidth := 0
	for _, m := range *matches {
		setMax(len(m.file), &fileWidth)
		setMax(digitCount(m.line), &lineWidth)
		setMax(len(m.version.format(format)), &versionWidth)
	}

	return fileWidth, lineWidth, versionWidth
}

func IndexOf[T string | int](list *[]T, value T) int {
	index := -1
	for i, listValue := range *list {
		if listValue == value {
			index = i
		}
	}
	return index
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func splitFileAndLineNotation(filePath string) (string, string) {
	parts := strings.Split(filePath, ":")
	if len(parts) == 1 {
		return filePath, ""
	}
	return parts[0], parts[1]
}
