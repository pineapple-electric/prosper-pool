// +build darwin

package main

import (
	"path/filepath"
)

const bundleIdentifier = "io.prosperpool.ProsperMiner"

func getSystemConfigFilePath() (string, error) {

	configFile := "prosper-miner.toml"

	path := filepath.Join("/Library", "Application Support", bundleIdentifier, configFile)
	return path, nil
}
