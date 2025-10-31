/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package plugins

import (
	"embed"

	"github.com/eolinker/apinto-dashboard/pm3/embed_registry"
)

var (
	//go:embed embed
	pluginDir embed.FS
)

func init() {
	err := embed_registry.LoadPlugins(&pluginDir, "embed", "plugin.yml")
	if err != nil {
		panic(err)
	}
}
