// +build freebsd linux netbsd openbsd

package main

func getSystemConfigFilePath() (string, error) {
	return "/etc/prosper-pool/prosper-pool.toml", nil
}
