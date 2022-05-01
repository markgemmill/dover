package app

import (
	"fmt"
	"github.com/elliotchance/orderedmap/v2"
	c "github.com/fatih/color"
	"github.com/marco-m/docopt-go"
	"os"
	"strings"
)

const VERSION = "0.2.1-dev.2"

func selectFormat(args ExecutionArgs, cfg ConfigValues) string {
	if args.format != "" {
		return args.format
	}
	return cfg.format
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
	initialize bool
	increment  bool
	format     string
	verbose    bool
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

func (u *Usage) writeAppPrefix(b *strings.Builder, writer *ColorizedWriter, first bool) {
	prefix := u.name
	if !first {
		prefix = fmt.Sprintf("%-0*s", len(u.name), " ")
	}
	writer.tool.Fprintf(b, "  %s", prefix)
}

func (u *Usage) writeCommandUsage(b *strings.Builder, writer *ColorizedWriter, cmdWidth int, cmd string, use string) {
	fmt.Fprintf(b, " %-0*s%s\n", cmdWidth, cmd, use)
}

func (u *Usage) buildUsage(b *strings.Builder, writer *ColorizedWriter) {
	writer.header.Fprint(b, "Usage:\n")

	for _, cmd := range u.usage.Keys() {
		usageList, _ := u.usage.Get(cmd)
		if len(usageList) > 0 {
			for index, use := range usageList {
				u.writeAppPrefix(b, writer, index == 0)
				u.writeCommandUsage(b, writer, len(cmd), cmd, use)
			}
		} else {
			u.writeAppPrefix(b, writer, true)
			u.writeCommandUsage(b, writer, len(cmd), cmd, "")
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
		name: "dover",
		description: `(do version) reports and updates your version number.
https://github.com/markgemmill/dover`,
		usage:   orderedmap.NewOrderedMap[string, []string](),
		options: orderedmap.NewOrderedMap[string, string](),
	}
	usageBuilder.addUsage("", []string{
		"[--increment] [--format=<fmt>] [--verbose]",
		"[--major | --minor | --patch | --build] ",
		"[--pre-release | --dev | --alpha | --beta | --rc | --release]",
	})
	usageBuilder.addUsage("init", []string{})

	usageBuilder.addOption("-i --increment", "Apply the increment.")
	usageBuilder.addOption("-f --format=<fmt>", "Apply format string: 000[-.+][(aA)[-.]0]")
	usageBuilder.addOption("-M --major", "Update major version segment.")
	usageBuilder.addOption("-m --minor", "Update minor version segment.")
	usageBuilder.addOption("-p --patch", "Update patch version segment.")
	usageBuilder.addOption("-P --pre-release", "Update to next pre-release.")
	usageBuilder.addOption("-d --dev", "Update dev version segment or bump dev build.")
	usageBuilder.addOption("-a --alpha", "Update alpha pre-release segment or bump alpha build.")
	usageBuilder.addOption("-b --beta", "Update beta pre-release segment or bump beta build.")
	usageBuilder.addOption("-r --rc", "Update release candidate segment or bump rc build.")
	usageBuilder.addOption("-B --build", "Update the pre-release build number.")
	usageBuilder.addOption("-R --release", "Clear pre-release version.")
	usageBuilder.addOption("-v --verbose", "Display details when incrementing.")
	usageBuilder.addOption("-h --help", "Display this help message.")
	usageBuilder.addOption("--version", "Display dover version.")

	return &usageBuilder
}

func ParseCommandline() docopt.Opts {
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
	return arguments

}

func compileArguments(opts docopt.Opts) ExecutionArgs {
	initialize, _ := opts.Bool("init")
	increment, _ := opts.Bool("--increment")
	format, _ := opts.String("--format")
	verbose, _ := opts.Bool("--verbose")

	args := ExecutionArgs{
		initialize: initialize,
		increment:  increment,
		format:     format,
		verbose:    verbose,
		part:       filterFlags(opts, []string{"major", "minor", "patch", "build"}),
		preRelease: filterFlags(opts, []string{"pre-release", "dev", "alpha", "beta", "rc", "release"}),
	}
	return args
}

func Execute() {
	args := compileArguments(ParseCommandline())

	c.NoColor = false

	cfg, err := configValues()
	ExitOnError(err)

	args.format = selectFormat(args, cfg)
	allMatches := getAllVersionStringMatches(cfg.files)

	if args.initialize == true {
		initialize(args)
		return
	}

	if !args.increment && args.part == "" && args.preRelease == "" {
		displayCurrentVersion(args, allMatches)
		return
	}

	if !args.increment && (args.part != "" || args.preRelease != "") {
		displayNextVersion(args, allMatches)
		return
	}

	if args.increment && (args.part != "" || args.preRelease != "") {
		applyNextVersion(args, allMatches)
		return
	}

}
