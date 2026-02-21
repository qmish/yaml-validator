package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

const watchDebounce = 300 * time.Millisecond

// runWatch запускает валидацию при изменении файлов.
func runWatch(files []string) {
	baseCfg := loadConfig()
	runValidation(baseCfg, files)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "watch: %v\n", err)
		os.Exit(1)
	}
	defer watcher.Close()

	for _, f := range files {
		absPath, err := filepath.Abs(f)
		if err != nil {
			absPath = f
		}
		if err := watcher.Add(absPath); err != nil {
			// файл мог быть удалён; продолжаем с остальными
			continue
		}
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	debounce := time.NewTimer(0)
	if !debounce.Stop() {
		<-debounce.C
	}
	defer debounce.Stop()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				debounce.Reset(watchDebounce)
			}
		case <-debounce.C:
			fmt.Println("\n--- change detected ---")
			runValidation(baseCfg, files)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Fprintf(os.Stderr, "watch error: %v\n", err)
		case sig := <-sigCh:
			fmt.Printf("\nwatch: %v, exiting\n", sig)
			return
		}
	}
}
