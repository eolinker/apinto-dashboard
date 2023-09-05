//go:build !windows

package local

import (
	"path/filepath"
)

var PluginDir, _ = filepath.Abs("local")

func cmdPath(name string) string {
	return filepath.Join(PluginDir, name)
}
