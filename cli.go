package main

import (
	"fmt"
	"github.com/elliotchance/orderedmap/v2"
	c "github.com/fatih/color"
	a "github.com/logrusorgru/aurora"
	"github.com/marco-m/docopt-go"
	"os"
	"strings"
)

func setMax(currentValue int, maxValue *int) {
	if currentValue > *maxValue {
		*maxValue = currentValue
	}
}

func maxLength(items []string) int {
	var maxLen int = 0
	for _, item := range items {
		itemLen := len(item)
		if itemLen > maxLen {
			maxLen = itemLen
		}
	}
	return maxLen
}

func keys(mapping map[string]string) []string {
	var list = []string{}
	for k, _ := range mapping {
		list = append(list, k)
	}
	return list
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
	var fileWidth int = 0
	var lineWidth int = 0
	var versionWidth int = 0
	for _, m := range *matches {
		setMax(len(m.file), &fileWidth)
		setMax(digitCount(m.line), &lineWidth)
		setMax(len(m.version.format(format)), &versionWidth)
	}

	return fileWidth, lineWidth, versionWidth
}

func printVersionChanges(matches *[]*VersionMatch, part string, release string, format string, updated bool) {
	var fileW, lineW, versW int
	fileW, lineW, versW = getMaxColumnWidths(matches, format)

	_update := ""
	if updated == true {
		_update = "updated "
	}

	for _, match := range *matches {
		nv := match.version.bump(part, release)
		fmt.Printf("%-0*s: %0*d %s%-0*s -> %s\n", fileW, a.Yellow(match.file), lineW, a.Blue(match.line), _update, versW, a.BrightWhite(match.version.format(format)).Bold(), a.BrightWhite(nv.format(format)))
	}
}

func printInconsistentVersions(matches *[]*VersionMatch, format string) {
	var fileW, lineW, versW int
	fileW, lineW, versW = getMaxColumnWidths(matches, format)

	for _, match := range *matches {
		fmt.Printf("%-0*s: %0*d  %-0*s\n", fileW, a.Yellow(match.file), lineW, a.Blue(match.line), versW, a.BrightWhite(match.version.format(format)).Bold())
	}
}

func selectFormat(args ExecutionArgs, cfg ConfigValues) string {
	if args.format != "" {
		return args.format
	}
	return cfg.format
}

func doVersion(args ExecutionArgs, matches *[]*VersionMatch) {
	if !assertVersionMatchConsistency(matches) {
		fmt.Print(a.BrightMagenta("\nVersions do not match across all files.\n"))
		printInconsistentVersions(matches, args.format)
	} else {
		fmt.Println((*matches)[0].version.format(args.format))
	}
}

func doNext(args ExecutionArgs, matches *[]*VersionMatch) {
	if !assertVersionMatchConsistency(matches) {
		fmt.Print(a.BrightMagenta("Versions do not match across all files.\n"))
		printInconsistentVersions(matches, args.format)
	} else {
		printVersionChanges(matches, args.part, args.preRelease, args.format, false)
	}
}

func doBump(args ExecutionArgs, matches *[]*VersionMatch) {
	if !assertVersionMatchConsistency(matches) {
		fmt.Print(a.BrightMagenta("No files have been changed!\n"))
		fmt.Print(a.BrightMagenta("Versions do not match across all files.\n"))
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

func filterFlags(args map[string]any, flags []string) string {
	activeFlags := []string{}
	for key, value := range args {
		key = strings.TrimLeft(key, "-")
		for _, flag := range flags {
			if key == flag && value == true {
				activeFlags = append(activeFlags, key)
			}
		}
	}
	switch len(activeFlags) {
	case 1:
		return activeFlags[0]
	case 0:
		return ""
	default:
		panic("Active flags can only be 1 or 0.")
	}
}

type ExecutionArgs struct {
	increment  bool
	format     string
	part       string
	preRelease string
}

type ColorizedWriter struct {
	tool   *c.Color
	header *c.Color
	flags  *c.Color
	desc   *c.Color
}

type Usage struct {
	name        string
	description string
	usage       *orderedmap.OrderedMap[string, []string]
	options     *orderedmap.OrderedMap[string, string]
}

func (u *Usage) addUsage(cmd string, args []string) {
	u.usage.Set(cmd, args)
}

func (u *Usage) addOption(optFlag string, optText string) {
	u.options.Set(optFlag, optText)
}

func (u *Usage) buildUsage(b *strings.Builder, writer *ColorizedWriter) {
	writer.header.Fprint(b, "Usage:\n")

	cmdWidth := maxLength(u.usage.Keys())

	if cmdWidth > 0 {
		cmdWidth = cmdWidth + 1
	}

	for _, cmd := range u.usage.Keys() {
		usageList, _ := u.usage.Get(cmd)
		for index, use := range usageList {
			prefix := u.name
			if index != 0 {
				prefix = fmt.Sprintf("%-0*s", len(u.name), " ")
			}
			writer.tool.Fprintf(b, "  %s", prefix)
			fmt.Fprintf(b, " %-0*s%s\n", cmdWidth, cmd, use)
		}
	}

	writer.tool.Fprintf(b, "  %s", u.name)
	fmt.Fprintf(b, " --help\n")

	writer.tool.Fprintf(b, "  %s", u.name)
	fmt.Fprintf(b, " --version\n")
}

func (u *Usage) buildOptions(b *strings.Builder, writer *ColorizedWriter) {

	// OPTIONS
	fmt.Fprint(b, "\n")
	writer.header.Fprint(b, "Options:\n")

	optionsWidth := maxLength(u.options.Keys())

	for _, opt := range u.options.Keys() {
		optDesc, _ := u.options.Get(opt)
		writer.flags.Fprintf(b, "  %-0*s", optionsWidth, opt)
		writer.desc.Fprintf(b, "  %s\n", optDesc)
	}
}

func (u *Usage) HelpText(colorize bool) string {
	c.NoColor = !colorize
	writer := ColorizedWriter{
		tool:   c.New(c.FgYellow),
		header: c.New(c.FgRed, c.Underline),
		flags:  c.New(c.FgHiWhite),
		desc:   c.New(c.FgBlue),
	}

	var b strings.Builder

	writer.tool.Fprintf(&b, "\n%s", u.name)
	fmt.Fprintf(&b, " %s\n\n", u.description)

	// USAGE
	u.buildUsage(&b, &writer)
	u.buildOptions(&b, &writer)

	return b.String()
}

func (u *Usage) UsageText(colorize bool) string {
	c.NoColor = !colorize
	writer := ColorizedWriter{
		tool:   c.New(c.FgYellow),
		header: c.New(c.FgRed, c.Underline),
		flags:  c.New(c.FgHiWhite),
		desc:   c.New(c.FgBlue),
	}

	var b strings.Builder

	writer.tool.Fprintf(&b, "\n%s\n\n", "Invalid arguments provided...")

	// USAGE
	u.buildUsage(&b, &writer)

	return b.String()
}

func NewUsageBuilder() *Usage {

	usageBuilder := Usage{
		name:        "dover",
		description: `(do version) reports and updates your version number.`,
		usage:       orderedmap.NewOrderedMap[string, []string](),
		options:     orderedmap.NewOrderedMap[string, string](),
	}
	usageBuilder.addUsage("", []string{
		"[--increment] [--format=<fmt>] ",
		"[--major | --minor | --patch | --build] ",
		"[--dev | --alpha | --beta | --rc | --release]",
	})
	usageBuilder.addOption("-i --increment", "Apply the increment.")
	usageBuilder.addOption("-f --format=<fmt>", "Apply format string.")
	usageBuilder.addOption("-M --major", "Update major version segment.")
	usageBuilder.addOption("-m --minor", "Update minor version segment.")
	usageBuilder.addOption("-p --patch", "Update patch version segment.")
	usageBuilder.addOption("-d --dev", "Update dev version segment.")
	usageBuilder.addOption("-a --alpha", "Update alpha pre-release segment.")
	usageBuilder.addOption("-b --beta", "Update beta pre-release segment.")
	usageBuilder.addOption("-r --rc", "Update release candidate segment.")
	usageBuilder.addOption("-B --build", "Update the pre-release build number.")
	usageBuilder.addOption("-R --release", "Clear pre-release version.")
	usageBuilder.addOption("-h --help", "Display this help message")
	usageBuilder.addOption("--version", "Display dover version.")

	return &usageBuilder
}

func ParseCommandline() ExecutionArgs {
	var version = fmt.Sprintf("dover v%s", VERSION)
	usage := NewUsageBuilder()

	var PrintHelpAndExit = func(err error, output string) {
		/// output could be --version or --help

		outputUsage := strings.Contains(output, "Usage:")
		outputOptions := strings.Contains(output, "Options:")

		if outputUsage == true && outputOptions == true {
			// if --help
			output = usage.HelpText(true)
		} else if outputUsage == true && outputOptions == false {
			// if invalid arguemnts
			output = usage.UsageText(true)
		}
		// else --version - use default output

		if err != nil {
			fmt.Fprintln(os.Stderr, output)
			os.Exit(1)

		} else {
			fmt.Println(output)
			os.Exit(0)
		}
	}

	var CLIParser = &docopt.Parser{
		HelpHandler:   PrintHelpAndExit,
		OptionsFirst:  false,
		SkipHelpFlags: false,
	}

	arguments, _ := CLIParser.ParseArgs(usage.HelpText(false), nil, version)

	increment, _ := arguments.Bool("--increment")
	format, _ := arguments.String("--format")

	args := ExecutionArgs{
		increment:  increment,
		format:     format,
		part:       filterFlags(arguments, []string{"major", "minor", "patch", "build"}),
		preRelease: filterFlags(arguments, []string{"dev", "alpha", "beta", "rc", "release"}),
	}

	return args
}

func Execute() {
	args := ParseCommandline()

	cfg, err := configValues()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	args.format = selectFormat(args, cfg)
	allMatches := getAllVersionStringMatches(cfg.files)

	if !args.increment && args.part == "" && args.preRelease == "" {
		doVersion(args, allMatches)
		return
	}

	if !args.increment && (args.part != "" || args.preRelease != "") {
		doNext(args, allMatches)
		return
	}

	if args.increment && (args.part != "" || args.preRelease != "") {
		doBump(args, allMatches)
		return
	}

}
