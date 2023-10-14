package pkg

import (
	"fmt"

	"github.com/barelyhuman/go/poller"
)

type Watcher struct {
	poller        *poller.Poller
	pathsToWatch  []string
	pathsToIgnore []string
}

func (w *Watcher) Wait() chan struct{} {
	return w.poller.Wait()
}

func (w *Watcher) OnEvent(exec func(poller.PollerEvent)) {
	w.poller.OnEvent(exec)
}

func (w *Watcher) Stop() {
	w.poller.Stop()
}

func (w *Watcher) Start() {
	for _, toIgnore := range w.pathsToIgnore {
		w.poller.IgnorePattern(toIgnore)
	}
	for _, p := range w.pathsToWatch {
		w.poller.Add(p)
	}
	w.poller.Start()

	w.poller.OnEvent(func(e poller.PollerEvent) {
		fmt.Printf("event: %v\n", e.Path)
	})
}

type WatcherModule func(*Watcher)

func CreateNewWatcher(mods ...WatcherModule) (w *Watcher, err error) {
	w = &Watcher{}
	w.poller = poller.NewPollWatcher(60)
	for _, mod := range mods {
		mod(w)
	}
	return
}

func WatcherWithPath(pathToWatch ...string) WatcherModule {
	return func(w *Watcher) {
		w.pathsToWatch = append(w.pathsToWatch, pathToWatch...)
	}
}

func WatcherWithIgnoredPath(pathsToIgnore ...string) WatcherModule {
	return func(w *Watcher) {
		w.pathsToIgnore = append(w.pathsToIgnore, pathsToIgnore...)
	}
}
