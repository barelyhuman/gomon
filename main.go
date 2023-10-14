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
		Usage:           "A go program executor with a file watcher",
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
					&cli.IntFlag{
						Name:    "poll",
						Aliases: []string{"p"},
						Value:   1000,
						Usage:   "duration of `POLLING` time in milliseconds (eg: 2 seconds would be 2000)",
					},
					&cli.StringSliceFlag{
						Name:    "include",
						Aliases: []string{"i"},
						Usage:   "`PATH`(s) to include, you can add multiple -i flags to add more paths ",
					},
					&cli.StringSliceFlag{
						Name:    "exclude",
						Aliases: []string{"e"},
						Usage:   "`PATH`(s) to exclude, you can add multiple -e flags to add more paths ",
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
