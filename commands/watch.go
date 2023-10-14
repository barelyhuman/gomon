package commands

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"syscall"

	"github.com/barelyhuman/go/poller"
	"github.com/barelyhuman/gomon/pkg"
	"github.com/urfave/cli/v2"
)

func Watch(c *cli.Context) (err error) {
	pollingDuration := c.Int("poll")
	pathsToInclude := c.StringSlice("include")
	pathsToExclude := c.StringSlice("exclude")

	gOptions := []pkg.WatcherModule{}

	if len(pathsToInclude) > 0 {
		gOptions = append(gOptions, pkg.WatcherWithPath(pathsToInclude...))
	}

	if len(pathsToExclude) > 0 {
		gOptions = append(gOptions, pkg.WatcherWithIgnoredPath(pathsToExclude...))
	}

	watcher, err := pkg.CreateNewWatcher(pollingDuration, gOptions...)

	if err != nil {
		return err
	}

	watcher.Start()

	if c.NArg() > 1 {
		c.Args()
	}

	baseArg := c.Args().First()
	fmt.Printf("baseArg: %vn", baseArg)
	pgid := runCmd(baseArg)
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan
		syscall.Kill(-pgid, syscall.SIGKILL)
		watcher.Stop()
	}()

	watcher.OnEvent(func(pe poller.PollerEvent) {
		syscall.Kill(-pgid, syscall.SIGKILL)
		pgid = runCmd(baseArg)
	})

	<-watcher.Wait()

	return nil
}

func runCmd(path string) (pgid int) {
	cmd := exec.Command("go", "run", path)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Start()
	id, _ := syscall.Getpgid(cmd.Process.Pid)
	return id
}
