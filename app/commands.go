package app

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

func printCurrentVersions(matches *[]*VersionMatch, format string) {
	var fileW, lineW, versW int
	fileW, lineW, versW = getMaxColumnWidths(matches, format)

	for _, match := range *matches {
		fmt.Printf(
			"%-0*s: %0*d  %-0*s\n",
			fileW,
			aurora.Yellow(match.file),
			lineW,
			aurora.Blue(match.line),
			versW,
			aurora.BrightWhite(match.version.format(format)).Bold(),
		)
	}
}

func printVersionChanges(matches *[]*VersionMatch, part string, release string, format string, updated bool) {
	var fileW, lineW, versW int
	fileW, lineW, versW = getMaxColumnWidths(matches, format)

	_update := ""
	if updated {
		_update = "updated "
	}

	for _, match := range *matches {
		nv := match.version.bump(part, release)
		fmt.Printf(
			"%-0*s: %0*d %s%-0*s -> %s\n",
			fileW,
			aurora.Yellow(match.file),
			lineW,
			aurora.Blue(match.line),
			_update,
			versW,
			aurora.BrightWhite(match.version.format(format)).Bold(),
			aurora.BrightWhite(nv.format(format)),
		)
	}
}

func displayInconsistentVersionMatch(args ExecutionArgs, matches *[]*VersionMatch) {
	if !assertVersionMatchConsistency(matches) {
		if args.increment {
			fmt.Print(aurora.BrightMagenta("No files have been changed!\n"))
		}
		fmt.Print(aurora.BrightMagenta("\nVersions do not match across all files.\n"))
		printCurrentVersions(matches, args.format)
		os.Exit(1)
	}
}

func displayCurrentVersion(args ExecutionArgs, matches *[]*VersionMatch) {
	displayInconsistentVersionMatch(args, matches)

	if args.verbose {
		printCurrentVersions(matches, args.format)
		return
	}
	fmt.Println((*matches)[0].version.format(args.format))
}

func displayFutureVersion(args ExecutionArgs, matches *[]*VersionMatch) {
	displayInconsistentVersionMatch(args, matches)

	newVers := (*matches)[0].version.bump(args.part, args.preRelease)
	version := newVers.format(args.format)
	fmt.Println(version)
}

func displayNextVersion(args ExecutionArgs, matches *[]*VersionMatch) {
	displayInconsistentVersionMatch(args, matches)

	printVersionChanges(matches, args.part, args.preRelease, args.format, false)
}

func applyNextVersion(args ExecutionArgs, matches *[]*VersionMatch) {
	displayInconsistentVersionMatch(args, matches)

	var newVers Version

	for _, match := range *matches {
		newVers = match.version.bump(args.part, args.preRelease)
		version := newVers.format(args.format)
		writeVersionUpdate(match.file, match.line, version)
	}

	if args.verbose {
		printVersionChanges(matches, args.part, args.preRelease, args.format, true)
	} else {
		fmt.Println(newVers.format(args.format))
	}
}

func initialize() {
	if fileExists(DOVER_CONFIG_FILE) {
		fmt.Println(aurora.BrightMagenta("Dover configuration file `.dover` already exists!"))
		return
	}

	err := os.WriteFile(DOVER_CONFIG_FILE, []byte(DOVER_DEFAULT_CONFIG), 0666)
	ExitOnError(err)
	fmt.Println(aurora.BrightGreen("Default `.dover` configuration file created."))
	fmt.Println("*** Be sure to add your project's versioned files! ***")
}
