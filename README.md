![dover-logo](./dover-logo.svg)

[![Go Report Card](https://goreportcard.com/badge/github.com/markgemmill/dover)](https://goreportcard.com/report/github.com/markgemmill/dover)

*version 0.2.1-dev.1*

A commandline utility for incrementing your project version numbers.

## What does it do?

When `dover` is run from the root directory of your project, it does the following:

1. looks for a configuration file (.dover, pyproject.toml, package.json)

2. reads the first available dover configuration file/format:

     `.dover` - dover toml format:

               [dover]
               version_format = "000.a0"
               versioned_files = [
                   "main.go"
               ]

     `pyproject.toml` - python project file:
 
        [tool.dover]
        version_format = "000.r0"
        versioned_files = ["pyproject.toml", "dover/cli.py"]

     `package.json` - javascript package file:

       {"name": "package-name",
        "version": "0.0.0",
        "dover": {
                "version_format": "000.r0"
                "versioned_files": ["package.json"] 
       }}



3. searches “version” strings in the files listed under `versioned_files`
   1. file paths can be appended with a list of line numbers (e.g.  "README.md:2,10") to 
      restrict which lines searches for version numbers.
4. validates all version strings are the same across all files.
5. performs the following based on arguments:

    `dover` - without args will display the projects current version number.
 
    `dover -M` - with a modifying flag and **WITHOUT** the -i option: displays the files that would be updated.
 
    `dover -iM` - with a modifying flag and the -i option: **APPLIES** the update and displays the new version number.
 


## Usage

    ... dover --help

    dover (do version) reports and updates your version number.

    Usage:
      dover [--increment] [--format=<fmt>] [--verbose]
            [--major | --minor | --patch | --build]
            [--dev | --alpha | --beta | --rc | --release]
      dover init
      dover --help
      dover --version

    Options:
      -i --increment     Apply the increment.
      -f --format=<fmt>  Apply format string.
      -M --major         Update major version segment.
      -m --minor         Update minor version segment.
      -p --patch         Update patch version segment.
      -d --dev           Update dev version segment.
      -a --alpha         Update alpha pre-release segment.
      -b --beta          Update beta pre-release segment.
      -r --rc            Update release candidate segment.
      -B --build         Update the pre-release build number.
      -R --release       Clear pre-release version.
      -v --verbose       Display details when incrementing.
      -h --help          Display this help message
      --version          Display dover version.



### Fetch Current Project Version

Calling dover without any arguments returns the current version number of the project:

    ... dover
    0.1.0-alpha.0

Using the format flag will format the version accordingly:

    ... dover -f 000r0
    0.1.0a0

If there are multiple versioned files and the versions are out of sync:
   
   ... dover
   README.md:2 0.1.0-dev.0
   main.go  :3 0.2.0

### Reviewing Version Increment Changes

Calling dover with one the segment options (e.g. --minor), will print a listing of the propsed version change and the files that will be effected:

    ... dover --minor
    setup.py      10 0.1.0 -> 0.2.0
    setup.cfg     02 0.1.0 -> 0.2.0
    dover/cli.py  13 0.1.0 -> 0.2.0

Attention:
    Only the use of the `–i, --increment` option will perform an update to your files.

### Applying Version Increment Changes

To save the change make the same call with the `-i, --increment` option:

    ... dover --minor --increment
    0.2.0

or 

    ... dover -mi
    0.2.0

Using the `-v, --verbose` flag will show all file changes:

    ... dover --miv
    setup.py      10 0.1.0 -> 0.2.0
    setup.cfg     02 0.1.0 -> 0.2.0
    dover/cli.py  13 0.1.0 -> 0.2.0


### Pre-Release Options

Applying a pre-release option (–dev, –alpha, –beta or –rc) appends the pre-release to the current version:

    ... dover --dev
    setup.py      10 0.1.0 -> 0.1.0-dev.0
    setup.cfg     02 0.1.0 -> 0.1.0-dev.0
    dover/cli.py  13 0.1.0 -> 0.1.0-dev.0


Subsequent call to the current pre-release increments the build version:

    ... dover --dev
    setup.py    : 10 0.1.0-dev.0 -> 0.1.0-dev.1
    setup.cfg   : 02 0.1.0-dev.0 -> 0.1.0-dev.1
    dover/cli.py: 13 0.1.0-dev.0 -> 0.1.0-dev.1


Applying the `-B, --build` option increments the current pre-release build version as well:

    ... dover -B
    setup.py      10 0.1.0-dev.0 -> 0.1.0-dev.1
    setup.cfg     02 0.1.0-dev.0 -> 0.1.0-dev.1
    dover/cli.py  13 0.1.0-dev.0 -> 0.1.0-dev.1


Use the `-R, --release` option to move from a pre-release version to the production version:

    ... dover -R
    setup.py      10 0.1.0-dev.0 -> 0.1.0
    setup.cfg     02 0.1.0-dev.0 -> 0.1.0
    dover/cli.py  13 0.1.0-dev.0 -> 0.1.0


## Version Formats

The default version format dover uses is:

    major.minor.patch[-(dev|alpha|beta|rc).version]

The output format can be controlled with the `-f, --format` option spec:

    000[(.|-|+)](r|R)[(.|-)]0

| Segment             | Format Value | Note                                                                                              |
|---------------------|--------------|---------------------------------------------------------------------------------------------------|
| Major.Minor.Patch	  | 000	         | Required.                                                                                      |
| Separator	          | . - +        | Optional. Dash, dot or plus.                                                                   |
| Pre-Release         | r *or* R     | Optional. Defaults to A.<br/>a = short name: d, a, b, rc <br/>A = long name: dev, alpha, beta, rc |
| Separator	          | . -	         | Optional. Dash or dot.                                                                         |
| Pre-Release Version | 0            | Version will always display if there is a pre-release.                                            |


Format examples:

| dover command    | output      | Note                 |
|------------------|-------------|----------------------|
| `dover           ` | 0.4.0-dev.1 | *default format*   |
| `dover –f 000-R.0` | 0.4.0-dev.1 | *default format*   |
| `dover –f 000+R.0` | 0.4.0+dev.1 |                    |
| `dover –f 000.R.0` | 0.4.0.dev.1 |                    |
| `dover –f 000R   ` | 0.4.0dev1   |                    |
| `dover –f 000r   ` | 0.4.0d1     |                    |
| `dover –f 000.r  ` | 0.4.d1      |                    |
| `dover –f 000r0  ` | 0.4d1       |                    |
| `dover –f 000-r0 ` | 0.4-d1      |                    |
| `dover –f 000-r-0` | 0.4-d-1     |                    |


### What If There Is a Problem?

If at any point the version numbers between multiple files being tracked are miss-aligned, dover will raise an error:

    ... dover -iB
    Versions do not match across all files.
    package.json: 2  0.1.1-alpha.0
    main.go     : 2  0.1.1-alpha.2
