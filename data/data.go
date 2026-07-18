package data

import (
	"os"
	"path/filepath"
)

func getDataDir() string {
	base := os.Getenv("APPDATA")
	if base == "" {
		base = filepath.Join(os.Getenv("HOME"), ".local", "share")
	}
	return filepath.Join(base, "OpenCli", "persistent")
}
