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
	watchPathsFlag := flag.String("w", ".", "`PATHS` to watch, separated by commas(,)")
	ignorePathsFlag := flag.String("exclude", "", "`PATHS` to exclude, separated by commas(,)")
	recursePathFlag := flag.Bool("r", true, "watch the mentioned folders recursively")

	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("Please provide a file to run on change")
		os.Exit(1)
	}

	execPath := (flag.Args())[0]

	paths := strings.Split(*watchPathsFlag, ",")

	// we append the glob paths with existing glob patterns because they won't
	// match with most actual file names and in case a pattern wasn't supposed to be glob
	// it doesn't get ignored accidentally since filepath.Glob will return empty if theres no files
	// matching the glob
	ignorePaths := trimPaths(strings.Split(*ignorePathsFlag, ","))
	ignorePaths = append(ignorePaths, normalizeGlobs(ignorePaths)...)

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

	var watchableDirs = paths[:]

	if *recursePathFlag {
		for _, walkPath := range paths {
			filepath.Walk(walkPath, func(subPath string, info os.FileInfo, err error) error {

				if pathInPathList(subPath, ignorePaths) {
					return nil
				}

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

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if pathInPathList(event.Name, ignorePaths) {
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

	<-done
}

func pathInPathList(toSearch string, list []string) bool {
	for _, toIgnore := range list {
		if strings.Contains(toSearch, toIgnore) || strings.HasPrefix(toSearch, toIgnore) {
			return true
		}
	}
	return false
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

func trimPaths(paths []string) []string {
	trimmedIgnorePaths := []string{}
	for _, i := range paths {
		trimmedIgnorePaths = append(trimmedIgnorePaths, strings.TrimSpace(i))
	}
	return trimmedIgnorePaths
}

func normalizeGlobs(patterns []string) []string {
	matchedPaths := []string{}
	for _, pattern := range patterns {
		matched, err := filepath.Glob(pattern)
		if err != nil {
			//digest
			log.Println(pattern, ": invalid pattern")
			continue
		}
		matchedPaths = append(matchedPaths, matched...)
	}
	return matchedPaths
}
