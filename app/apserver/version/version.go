package version

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

// These should be set via go build -ldflags -X 'xxxx'.
var Version = "unknown"
var goversion = "unknown"
var gitcommit = "unknown"
var buildtime = "unknown"
var builduser = "unknown"

var profileInfo []byte

func init() {
	buffer := &bytes.Buffer{}
	fmt.Fprintf(buffer, "apserver version: %s\n", Version)
	fmt.Fprintf(buffer, "Golang version: %s\n", goversion)
	fmt.Fprintf(buffer, "Git commit hash: %s\n", gitcommit)
	fmt.Fprintf(buffer, "Built on: %s\n", buildtime)
	fmt.Fprintf(buffer, "Built by: %s\n", builduser)
	profileInfo = buffer.Bytes()
}
func Handler(ctx *gin.Context) {

	ctx.Writer.Write(profileInfo)

}
func Build() *cli.Command {
	return &cli.Command{
		Name: "version",
		Action: func(context *cli.Context) error {
			fmt.Print(string(profileInfo))
			return nil
		},
	}
}

func GetVersion() string {
	return Version
}
