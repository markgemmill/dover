package app

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestJSONConfigWithValues(t *testing.T) {
	projectFile := `{
	"name": "Some Project",
	"version": "0.0.0",
	"dover": {
		"version_format": "000.a0",
		"versioned_files": [
			"project.json"
		]
	}
}`
	cfg, err := parseJSONConfig(projectFile)

	assert.Nil(t, err)
	assert.Equal(t, "000.a0", cfg.format)
	assert.Equal(t, 1, len(cfg.files))
}

func TestJSONConfigWithValues2(t *testing.T) {
	projectFile := `{
	"name": "Some Project",
	"version": "",
	"dover": {
		"version_format": "000.a0",
		"versioned_files": [
			"package.json",
            "main.go"
		]
	}
}`
	cfg, err := parseJSONConfig(projectFile)

	assert.Nil(t, err)
	assert.Equal(t, "000.a0", cfg.format)
	assert.Equal(t, 2, len(cfg.files))
}

func TestJSONConfigWithOnlyVersionedFiles(t *testing.T) {
	projectFile := `{
	"name": "Some Project",
	"version": "0.0.0",
	"dover": {
		"versioned_files": [
			"project.json"
		]
	}
}`
	cfg, err := parseJSONConfig(projectFile)

	assert.Nil(t, err)
	assert.Equal(t, "", cfg.format)
	assert.Equal(t, 1, len(cfg.files))
}

func TestJSONConfigWithNoFiles(t *testing.T) {
	projectFile := `{
	"name": "Some Project",
	"version": "0.0.0",
	"dover": {
		"version_format": "0.0.0",
		"versioned_files": []
	}
}`
	cfg, err := parseJSONConfig(projectFile)

	assert.NotNil(t, err)
	assert.Equal(t, "", cfg.format)
	assert.Equal(t, 0, len(cfg.files))
}

func TestJSONConfigWithNoFileVar(t *testing.T) {
	projectFile := `{
	"name": "Some Project",
	"version": "0.0.0",
	"dover": {
		"version_format": "0.0.0"
	}
}`
	cfg, err := parseJSONConfig(projectFile)

	assert.NotNil(t, err)
	assert.Equal(t, "", cfg.format)
	assert.Equal(t, 0, len(cfg.files))
}

type ConfigTestSuite struct {
	suite.Suite
	homeDir string
	tempDir string
}

func (suite *ConfigTestSuite) writeFile(name, content string) {
	file := filepath.Join(suite.tempDir, name)
	err := os.WriteFile(file, []byte(content), 0666)
	ExitOnError(err)
}

func (suite *ConfigTestSuite) SetupTest() {
	suite.homeDir, _ = os.Getwd()
	// suite.tempDir, _ = ioutil.TempDir("", "gotest-*")
	suite.tempDir, _ = os.MkdirTemp("", "gotest-*")
	suite.writeFile("coding.go", `\nVERSION = "0.1.0-a0"\n`)
	suite.writeFile("overhill.go", `\n__version__ = "0.1.0-a0"\n`)
	err := os.Chdir(suite.tempDir)
	ExitOnError(err)
}

func (suite *ConfigTestSuite) TeardownTest() {
	err := os.Chdir(suite.homeDir)
	ExitOnError(err)
	err = os.RemoveAll(suite.tempDir)
	ExitOnError(err)
}

func (suite *ConfigTestSuite) TestNoConfigFiles() {
	_, err := configValues()
	suite.NotNil(err)
	suite.Equal("unable to find dover configuration", fmt.Sprint(err))
}

func (suite *ConfigTestSuite) TestInvalidDoverConfigFile() {
	suite.writeFile(".dover", `[dover]`)

	_, err := configValues()
	suite.NotNil(err)
	suite.Equal("`.dover` config has no versioned_files", fmt.Sprint(err))
}

func (suite *ConfigTestSuite) TestConfigWithInvalidVersionedFile() {
	suite.writeFile(".dover", `[dover]
versioned_files = [
	"dunnowherethisis.go",
	"overhill.go"
]
`)

	_, err := configValues()
	suite.NotNil(err)
	suite.Equal("no such file: dunnowherethisis.go", fmt.Sprint(err))
}

func (suite *ConfigTestSuite) TestValidDoverConfigFile() {
	suite.writeFile(".dover", `[dover]
versioned_files = [
	"coding.go",
	"overhill.go"
]
`)

	cfg, err := configValues()
	suite.Nil(err)
	suite.Equal("000.A.0", cfg.format)
	suite.Equal(2, len(cfg.files))
}

func (suite *ConfigTestSuite) TestValidPyProjectConfigFile() {
	// test the default format code...
	suite.writeFile("pyproject.toml", `[tool.dover]
versioned_files = [
	"coding.go",
	"overhill.go"
]
`)

	cfg, err := configValues()
	suite.Nil(err)
	suite.Equal("000.A.0", cfg.format)
	suite.Equal(2, len(cfg.files))
}

func (suite *ConfigTestSuite) TestValidPackageJsonConfigFile() {
	suite.writeFile("package.json", `{
	"name": "project",
	"version": "0.1.0.beta.0",
	"dover": {
		"version_format": "000+a0",
		"versioned_files": [
			"coding.go",
			"overhill.go"
		]
	}
}`)

	cfg, err := configValues()
	suite.Nil(err)
	suite.Equal("000+a0", cfg.format)
	suite.Equal(2, len(cfg.files))
}

func TestRunConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
