package cli

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

const (
	appName = "apserver-authorization"
)

type App struct {
	app *cli.App
}

func NewApp() *App {
	return &App{app: &cli.App{
		Name:     appName,
		Usage:    fmt.Sprintf("%s controller", appName),
		Commands: make([]*cli.Command, 0, 3),
	}}
}

func (a *App) AppendCommand(cmd ...*cli.Command) {
	a.app.Commands = append(a.app.Commands, cmd...)
}

func (a *App) Run(args []string) error {
	return a.app.Run(args)
}

func (a *App) Default() {
	a.AppendCommand(
		GenCert(),
		GenRSAKeys(),
	)
}
