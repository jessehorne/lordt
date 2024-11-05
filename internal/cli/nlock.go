package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bmatcuk/doublestar/v4"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/s3rj1k/go-fanotify/fanotify"
)

func nlockHelp() {
	fmt.Println(`Usage: lordt nlock <options...>

nlock (notify-based lock) is an experimental tool for preventing processes from opening/reading files while it is running.

Currently 'nlock' uses fanotify to get file open events and responds as soon as it can. This means that it isn't
very fast and the 'kill' option isn't reliable.'. Some tests showed processes being able to open and then read up to 
9kb before being killed.

---
Options
---
--pattern=...			- The file/directory pattern. (e.g 'lordt lock --pattern=/home/user/*.txt')
-log				- Print information on processes that try to access the locked files.
-kill				- Kill processes that attempt to access the files.
--conf=...			- Load a config file which replaces the need to provide these flags.

Example Usage:
lordt nlock --pattern=/home/user/**/*.txt -log -kill

Example config:

nlock.json
'''
{
	"patterns": [
		"/home/user/example.txt",
		"/home/user/**/*.docx"
	],
	"options": ["log", "kill"]
}
'''
`)
}

func NLockCommandHandler(ch *CommandHandler, args []string) error {
	if len(args) == 0 {
		nlockHelp()
		return nil
	}

	if args[0] == "help" {
		nlockHelp()
		return nil
	}

	confPath := flag.String("conf", "", "path to config file")
	pattern := flag.String("pattern", "", "pattern to match against file(s) that should be locked")
	shouldLog := flag.Bool("log", false, "should log information about processes that attempt to access locked file(s)")
	shouldKill := flag.Bool("kill", false, "should attempt to kill processes that attempt to access locked file(s)")
	flag.CommandLine.Parse(args)

	conf := &nlockConf{}

	if *confPath == "" {
		if *pattern == "" {
			log.Fatalln("pattern option required")
		}
		conf.Patterns = []string{*pattern}
		conf.ShouldLog = *shouldLog
		conf.ShouldKill = *shouldKill
	} else {
		c, err := loadConf(*confPath)
		if err != nil {
			log.Fatalln("couldn't load conf:", *confPath)
		}
		conf = c
	}

	notify, err := fanotify.Initialize(
		unix.FAN_CLOEXEC|
			unix.FAN_CLASS_NOTIF|
			unix.FAN_UNLIMITED_QUEUE|
			unix.FAN_UNLIMITED_MARKS,
		os.O_RDONLY|
			unix.O_LARGEFILE|
			unix.O_CLOEXEC,
	)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	var allFiles []string
	for _, p := range conf.Patterns {
		files, err := doublestar.FilepathGlob(p)
		if err != nil {
			continue
		}
		for _, f := range files {
			err := nlockFile(notify, f)
			if err != nil {
				log.Println(err)
				continue
			}
			path, err := filepath.Abs(f)
			if err != nil {
				continue
			}
			log.Println("locked:", path)
			allFiles = append(allFiles, path)
		}
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	f := func(notify *fanotify.NotifyFD) error {
		data, err := notify.GetEvent(os.Getpid())
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		if data == nil {
			return nil
		}

		defer data.Close()

		path, err := data.GetPath()
		if err != nil {
			return err
		}

		dataFile := data.File()
		defer dataFile.Close()

		fInfo, err := dataFile.Stat()
		if err != nil {
			return err
		}

		mTime := fInfo.ModTime()

		for _, p := range allFiles {
			if p == path {
				if data.MatchMask(unix.FAN_OPEN) {
					if conf.ShouldLog {
						log.Println(fmt.Sprintf("[OPEN] PID:%d %s - %v", data.GetPID(), path, mTime))
					}

					if conf.ShouldKill {
						err := syscall.Kill(data.GetPID(), syscall.SIGKILL)
						if err != nil {
							log.Println(err)
							continue
						}
						log.Println(fmt.Sprintf("[KILLED] PID:%d", data.GetPID()))
					}
				}
			}
		}

		return nil
	}

	for {
		select {
		case <-done:
			return nil
		default:
			err := f(notify)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return nil
}

type nlockConf struct {
	Patterns   []string
	ShouldLog  bool
	ShouldKill bool
}

type nlockConfJSON struct {
	Patterns []string
	Options  []string
}

func nlockFile(notify *fanotify.NotifyFD, path string) error {
	return notify.Mark(
		unix.FAN_MARK_ADD|unix.FAN_MARK_MOUNT,
		unix.FAN_MODIFY|unix.FAN_CLOSE_WRITE|unix.FAN_OPEN,
		unix.AT_FDCWD,
		path,
	)
}

func loadConf(confPath string) (*nlockConf, error) {
	data, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}

	var confJSON *nlockConfJSON
	err = json.Unmarshal(data, &confJSON)
	if err != nil {
		return nil, err
	}

	conf := &nlockConf{}
	conf.Patterns = confJSON.Patterns
	for _, opt := range confJSON.Options {
		if opt == "log" {
			conf.ShouldLog = true
		} else if opt == "kill" {
			conf.ShouldKill = true
		}
	}

	return conf, nil
}
