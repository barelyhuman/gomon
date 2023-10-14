package main

import (
	_ "embed"
	"fmt"
	"os"
	"time"

	"github.com/barelyhuman/gomon/commands"
	"github.com/urfave/cli/v2"
)

//go:embed .commitlog.release
var version string

func main() {
	app := &cli.App{
		Name:            "gomon",
		Usage:           "command executor with a file watcher",
		CommandNotFound: cli.ShowCommandCompletions,
		Action: func(c *cli.Context) error {
			return cli.ShowAppHelp(c)
		},
		Version:     version,
		Compiled:    time.Now(),
		HideVersion: false,
		Commands: []*cli.Command{
			{
				Name:    "watch",
				Aliases: []string{"w"},
				Usage:   "watch mode",
				Action: func(c *cli.Context) error {
					return commands.Watch(c)
				},
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:    "path",
						Aliases: []string{"p"},
						Usage:   "",
					},
					&cli.StringSliceFlag{
						Name:    "ignore",
						Aliases: []string{"i"},
						Usage:   "",
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[gomon] %v", err)
	}
}
