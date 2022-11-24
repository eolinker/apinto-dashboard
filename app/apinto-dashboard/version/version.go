package version

import (
	"bytes"
	"fmt"
)

// These should be set via go build -ldflags -X 'xxxx'.
var Version = "unknown"
var goVersion = "unknown"
var gitCommit = "unknown"
var buildTime = "unknown"
var buildUser = "unknown"
var eoscVersion = "unknown"

var profileInfo []byte

func init() {
	buffer := &bytes.Buffer{}
	fmt.Fprintf(buffer, "Apinto version: %s\n", Version)
	fmt.Fprintf(buffer, "Golang version: %s\n", goVersion)
	fmt.Fprintf(buffer, "Git commit hash: %s\n", gitCommit)
	fmt.Fprintf(buffer, "Built on: %s\n", buildTime)
	fmt.Fprintf(buffer, "Built by: %s\n", buildUser)
	fmt.Fprintf(buffer, "Built by eosc version: %s\n", eoscVersion)
	profileInfo = buffer.Bytes()
}

func Println() {
	fmt.Print(string(profileInfo))
}
