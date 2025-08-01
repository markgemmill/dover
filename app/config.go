package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/pelletier/go-toml"
)

func findConfigFile(fileName string) (string, error) {
	/*
		We're looking for dover config info in the following locations:

			.dover         (universal)
			pyproject.toml (python)
			package.json   (javascript)

	*/

	root := "./"
	fileSystem := os.DirFS(root)
	entries, err := fs.Glob(fileSystem, fileName)
	check(err)

	if len(entries) == 1 {
		return entries[0], nil
	}

	return "", fmt.Errorf("could not find %s config", fileName)
}

type ConfigValues struct {
	files  []string
	format string
}

type configParser func(string) (ConfigValues, error)

func getTomlConfigValues(configFile string) (ConfigValues, error) {
	/*
		Read the .dover configuration file
	*/
	cfg, err := toml.LoadFile(configFile)
	check(err)

	cfgV := ConfigValues{}

	getVersionedFiles := func(c *toml.Tree, pth string) []string {
		if c.Has(pth) {
			return c.GetArray(pth).([]string)
		}
		return []string{}
	}

	getVersionFormat := func(c *toml.Tree, pth string) string {
		if c.Has(pth) {
			return c.Get(pth).(string)
		}
		return ""
	}

	if cfg.Has("dover") {
		// .dover
		cfgV.files = getVersionedFiles(cfg, "dover.versioned_files")
		cfgV.format = getVersionFormat(cfg, "dover.version_format")
		return cfgV, nil
	} else if cfg.Has("tool.dover") {
		// pyproject.toml
		cfgV.files = getVersionedFiles(cfg, "tool.dover.versioned_files")
		cfgV.format = getVersionFormat(cfg, "tool.dover.version_format")
		return cfgV, nil
	}

	return cfgV, errors.New(fmt.Sprint("No dover config entries in ", configFile))
}

func readJSONConfig(configFile string) []byte {
	file, err := os.Open(configFile)
	check(err)
	defer file.Close()

	content, _ := io.ReadAll(file)
	return content
}

func parseJSONConfig(content string) (ConfigValues, error) {
	/*
		Read the project.json configuration file
	*/
	type ProjectJSON struct {
		Dover struct {
			VersionFormat  string   `json:"version_format"`
			VersionedFiles []string `json:"versioned_files"`
		} `json:"dover"`
	}

	cfgV := ConfigValues{}

	var payload ProjectJSON
	err := json.Unmarshal([]byte(content), &payload)
	if err != nil {
		return cfgV, fmt.Errorf("json parsing failed: %s", err)
	}

	if len(payload.Dover.VersionedFiles) == 0 {
		return cfgV, fmt.Errorf("no `dover` section or `dover.versioned_files` contains no file references")
	}

	cfgV.format = payload.Dover.VersionFormat
	cfgV.files = payload.Dover.VersionedFiles

	return cfgV, nil
}

func getJSONConfigValues(configFile string) (ConfigValues, error) {
	content := readJSONConfig(configFile)
	cfgV, err := parseJSONConfig(string(content))
	if err != nil {
		return cfgV, err
	}

	return cfgV, nil
}

const (
	DOVER_CONFIG_FILE        = ".dover"
	PYPROJECT_CONFIG_FILE    = "pyproject.toml"
	PACKAGE_JSON_CONFIG_FILE = "package.json"
)

const DOVER_DEFAULT_CONFIG = `[dover]
version_format = "000-A.0"
versioned_files = [
]
`

func configValues() (ConfigValues, error) {
	configOrder := []string{DOVER_CONFIG_FILE, PYPROJECT_CONFIG_FILE, PACKAGE_JSON_CONFIG_FILE}

	possibleConfigFiles := map[string]configParser{
		DOVER_CONFIG_FILE:        getTomlConfigValues,
		PYPROJECT_CONFIG_FILE:    getTomlConfigValues,
		PACKAGE_JSON_CONFIG_FILE: getJSONConfigValues,
	}

	var cfg ConfigValues

	for _, fileName := range configOrder {

		configParser := possibleConfigFiles[fileName]
		cfgFile, err := findConfigFile(fileName)
		if err != nil {
			continue
		}

		cfg, err := configParser(cfgFile)
		if err != nil {
			fmt.Printf("%s: %s", fileName, err)
			continue
		}

		if len(cfg.files) == 0 {
			return cfg, fmt.Errorf("`%s` config has no versioned_files", fileName)
		}

		for _, filePath := range cfg.files {
			filePath, _ := splitFileAndLineNotation(filePath)
			if !fileExists(filePath) {
				return cfg, fmt.Errorf("no such file: %s", filePath)
			}
		}

		if cfg.format == "" {
			cfg.format = "000.A.0"
		}

		return cfg, nil
	}

	return cfg, fmt.Errorf("unable to find dover configuration")
}
