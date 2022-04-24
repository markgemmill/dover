# dover v0.5.1

A commandline utility for incrementing your project version numbers.

## What does it do?

When `dover` is run from the root directory of your project, it does the following:

looks for a configuration file (.dover, pyproject.toml, package.json)

reads any dover configuration line in this format:

        [dover:file:relatvie/file.pth]

Or in the case of pyproject.toml:

    [tool.dover]
    versioned_files = ["pyproject.toml", "dover/cli.py"]

searches the configured file references for “version” strings
validates all version strings across all configured files.
displays and/or increments the version strings based upon cli options.


### Usage

... dover --help

    dover v0.5.1

    dover is a commandline utility for
    tracking and incrementing your
    project version numbers.

    Usage:
      dover [--list] [--debug] [--format=<fmt>]
      dover increment ((--major|--minor|--patch)
                       [--dev|--alpha|--beta|--rc] |
                       [--major|--minor|--patch]
                       (--dev|--alpha|--beta|--rc) | --release)
                       [--apply] [--debug] [--no-list] [--format=<fmt>]

    Options:
      -M --major      Update major version segment.
      -m --minor      Update minor version segment.
      -p --patch      Update patch version segment.
      -d --dev        Update dev version segment.
      -a --alpha      Update alpha pre-release segment.
      -b --beta       Update beta pre-release segment.
      -r --rc         Update release candidate segment.
      -R --release    Clear pre-release version.
      -x --no-list    Do not list files.
      --format=<fmt>  Apply format string.
      --debug         Print full exception info.
      -h --help       Display this help message
      --version       Display dover version.

### Basics

dover will look for either a .dover, setup.cfg, or pyproject.toml file within the current directory.

Files that are to be tracked by dover should be listed in the following format:

    [dover:file:relative-file-path.ext]

Example:

    [dover:file:setup.py]
    [dover:file:setup.cfg]
    [dover:file:dover/cli.py]

For projects using the pyproject.toml the format will be:

    [tool.dover]
    versioned_files = ["pyproject.toml", "dover/cli.py"]


#### Fetch Current Project Version

Calling dover without any arguments returns the current version number of the project.

    ... dover
    0.1.0

This is useful for capturing your project version for use with other commandline utilities.

### Currently Tracked File Status

Calling dover with the --list option, prints the project version and the list of all files and version strings being tracked, and the line they appear on:

    ... dover --list
    Current Version: 0.1.0
    Files:
        setup.py     0005 (__version__ = '0.1.0')
        setup.cfg    0002 (version = 0.1.0)
        dover/cli.py 0025 (__version__ = '0.1.0')

### Reviewing Version Increment Changes

Calling dover increment with one the the segment options (e.g. --minor), will print a listing of the propsed version change and the files that will be effected:

    ... dover increment --minor
    Current Version: 0.1.0
    New Version:     0.2.0
    Files:
        setup.py      (0.1.0 -> 0.2.0)
        setup.cfg     (0.1.0 -> 0.2.0)
        dover/cli.py  (0.1.0 -> 0.2.0)

Attention:
    Only the use of the –apply option will perform a update to your files.

### Applying Version Increment Changes

To save the change make the same call with the --apply option:

    ... dover increment --minor --apply
    Current Version: 0.1.0
    New Version:     0.2.0
    Files:
        setup.py      (0.1.0 -> 0.2.0)
        setup.cfg     (0.1.0 -> 0.2.0)
        dover/cli.py  (0.1.0 -> 0.2.0)
    Version updates applied.

### Pre-Release Options

Applying a pre-release option (–dev, –alpha, –beta or –rc), simply appends the pre-release to the current version:

    ... dover increment --alpha --no-list
    Current Version: 0.1.0
    New Version:     0.1.0-alpha

Tip:
    the –no-list option suppresses listing the files.

Applying a pre-release option to an existing pre-release of the same name increments the pre-release:

    ... dover increment --alpha --no-list
    Current Version: 0.1.0
    New Version:     0.1.0-alpha.1

Applying a pre-release option with a segment option, increments the segment and appends the pre-relase value:

    ... dover increment --minor --alpha --no-list
    Current Version: 0.1.0
    New Version:     0.2.0-alpha

Use the --release option to move from a pre-release version to the production version:

    ... dover increment  --release --no-list
    Current Version: 0.4.0-dev
    New Version:     0.4.0


## Version Formats

The default version format dover uses is:

    major.minor[.patch][-(dev|alpha|beta|rc)[.version]]

The output format can be controlled with the --format option with these options:

    000[(.|-|+)](a|A)[(.|-)]0

| Segment             | Format Value | Note                                                                                              |
|---------------------|--------------|---------------------------------------------------------------------------------------------------| 
| Major.Minor.Patch	  | 000	         | Required.                                                                                         | 
| Separator	          | . - +        | Optional. Dash, dot or plus.                                                                      |
| Pre-Release         | a *or* A     | Optional. Defaults to A.<br/>a = short name: d, a, b, rc <br/>A = long name: dev, alpha, beta, rc |
| Separator	          | . -	         | Optional. Dash or dot.                                                                            |
| Pre-Release Version | 0            | Version will always display if there is a pre-release.                                            |


Format examples:

| dover command    | output      | Note                 |
|------------------|-------------|----------------------|
| `dover           ` | 0.4.0-dev.1 | *default format*   |
| `dover –f 000-A.0` | 0.4.0-dev.1 | *default format*   |
| `dover –f 000+A.0` | 0.4.0+dev.1 |                    |
| `dover –f 000.A.0` | 0.4.0.dev.1 |                    |
| `dover –f 000A   ` | 0.4.0dev1   |                    |
| `dover –f 000a   ` | 0.4.0d1     |                    |
| `dover –f 000.a  ` | 0.4.d1      |                    |
| `dover –f 000a0  ` | 0.4d1       |                    |
| `dover –f 000-a0 ` | 0.4-d1      |                    |
| `dover –f 000-a-0` | 0.4-d-1     |                    |

## What If There Is a Problem?

If at any point the version numbers between files being tracked are miss-aligned, dover will raise an error:

    ... dover -i -M

Not all file versions match:

    setup.py      0.1.0  (__version__ = '0.1.0')
    setup.cfg     0.3.0  (version = 0.3.0)
    dover/cli.py  0.1.0  (__version__ = '0.1.0')