// +build linux

package main

func getDefaultHashTableDirectory() (string, error) {
	return "/var/lib/LXRHash"
}

func getSystemConfigFilePath() (string, error) {
	return "/etc/prosper-pool/prosper-pool.toml", nil
}
