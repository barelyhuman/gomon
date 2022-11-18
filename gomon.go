package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watchPathsFlag := flag.String("w", ".", "`PATHS` to watch")
	recursePathFlag := flag.Bool("r", true, "watch the mentioned folders recursively")

	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("Please provide a file to run on change")
		os.Exit(1)
	}

	execPath := (flag.Args())[0]

	paths := strings.Split(*watchPathsFlag, ",")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	defer watcher.Close()
	defer fmt.Println("")

	done := make(chan bool)
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt)

	var pgid int
	pgid = runCmd(execPath, true)

	go func() {
		<-signalChan
		syscall.Kill(-pgid, 15)
		done <- true
	}()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					syscall.Kill(-pgid, 15)
					pgid = runCmd(execPath, false)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("Error happened 😢", err)
			}
		}
	}()

	var watchableDirs = paths[:]

	if *recursePathFlag {
		for _, walkPath := range paths {
			filepath.Walk(walkPath, func(subPath string, info os.FileInfo, err error) error {
				if info.IsDir() {
					watchableDirs = append(watchableDirs, path.Join(walkPath, subPath))
				}
				return nil
			})
		}
	}

	for _, path := range watchableDirs {
		err = watcher.Add(path)

		if err != nil {
			panic(err)
		}
	}

	<-done
}

func runCmd(path string, first bool) (pgid int) {
	cmd := exec.Command("go", "run", path)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmd.Start()

	if first {
		formatPrint("Starting to watch")
	} else {
		formatPrint("Restarting due to change")
	}
	id, _ := syscall.Getpgid(cmd.Process.Pid)

	return id
}

func formatPrint(msg string) {
	fmt.Println("")
	fmt.Println("\x1b[36m*")
	fmt.Println("  " + msg)
	fmt.Println("*\x1b[0m")
	fmt.Println("")
}
