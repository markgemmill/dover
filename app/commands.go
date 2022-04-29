package app

import (
	"fmt"
	"github.com/logrusorgru/aurora"
)

func printInconsistentVersions(matches *[]*VersionMatch, format string) {
	var fileW, lineW, versW int
	fileW, lineW, versW = getMaxColumnWidths(matches, format)

	for _, match := range *matches {
		fmt.Printf("%-0*s: %0*d  %-0*s\n", fileW, aurora.Yellow(match.file), lineW, aurora.Blue(match.line), versW, aurora.BrightWhite(match.version.format(format)).Bold())
	}
}

func doVersion(args ExecutionArgs, matches *[]*VersionMatch) {
	if !assertVersionMatchConsistency(matches) {
		fmt.Print(aurora.BrightMagenta("\nVersions do not match across all files.\n"))
		printInconsistentVersions(matches, args.format)
	} else {
		fmt.Println((*matches)[0].version.format(args.format))
	}
}

func doNext(args ExecutionArgs, matches *[]*VersionMatch) {
	if !assertVersionMatchConsistency(matches) {
		fmt.Print(aurora.BrightMagenta("Versions do not match across all files.\n"))
		printInconsistentVersions(matches, args.format)
	} else {
		printVersionChanges(matches, args.part, args.preRelease, args.format, false)
	}
}

func doBump(args ExecutionArgs, matches *[]*VersionMatch) {
	if !assertVersionMatchConsistency(matches) {
		fmt.Print(aurora.BrightMagenta("No files have been changed!\n"))
		fmt.Print(aurora.BrightMagenta("Versions do not match across all files.\n"))
		printInconsistentVersions(matches, args.format)
	} else {
		for _, match := range *matches {
			newVers := match.version.bump(args.part, args.preRelease)
			version := newVers.format(args.format)
			writeVersionUpdate(match.file, match.line, version)
		}
		printVersionChanges(matches, args.part, args.preRelease, args.format, true)
	}
}
