package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJSONConfigWithValues(t *testing.T) {
	var projectFile = `{
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
	var projectFile = `{
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
	var projectFile = `{
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
	var projectFile = `{
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
	var projectFile = `{
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
