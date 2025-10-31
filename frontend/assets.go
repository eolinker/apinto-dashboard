package frontend

import "embed"

var (
	//go:embed dist/index.html
	indexContent []byte
	//go:embed dist
	dist embed.FS
	//go:embed logo/favicon.ico
	iconContent []byte

	//go:embed logo/logo.svg
	logoContent []byte
)
