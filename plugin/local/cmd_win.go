//go:build windows

package local

import (
	"fmt"
	"path/filepath"
)

var PluginDir, _ = filepath.Abs("local")

func cmdPath(name string) string {
	return filepath.Join(PluginDir, fmt.Sprintf("%s.exe", name))
}
