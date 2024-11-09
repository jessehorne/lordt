package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	honeyFileTypes = []string{"txt"}
)

func honeyHelp() {
	fmt.Println(`Usage: lordt honey <options...>

honey creates sets of decoy data at specified locations in your filesystem to hopefully catch ransomware before it exfiltrates or encrypts your real data.

---
Options
---
--path=...			- The directory where the decoy data will be created. (e.g '/home/user/decoy')
-b				- Breadth. How many files, roughly, per folder.
-d				- Depth. How deep we insert folders into folders. (e.g -d=3 would create a folder structure like /home/user/decoy/1/2/3/...)
-maxsize			- Max file size in KB
-minsize			- Min file size in KB
-filetypes			- Comma separated list of filetypes to generate (e.g docx,pdf,txt,sql,...)
--conf=...			- Load a config file which replaces the need to provide these flags.

-listfiletypes			- Print out the list of supported filetypes

Example Usage:
lordt honey --path=/home/user/random -b 10 -d 3 -maxsize 800 -minsize 9

Example config:

honey.json
'''
{
  "path": "./test/decoy",
  "b": 10,
  "d": 2,
  "maxsize": 9,
  "minsize": 1,
  "filetypes": ["pdf", "docx", "txt"]
}
'''
`)
}

func honeyListFileTypes() {
	fmt.Println(strings.Join(honeyFileTypes, ","))
}

type honeyConf struct {
	Path      string
	Breadth   int
	Depth     int
	MaxSize   int      // max generated file size in KB
	MinSize   int      // min generated file size in KB
	Filetypes []string // list of acceptable generated file types (e.g pdf, docx, ...)
}

func HoneyCommandHandler(ch *CommandHandler, args []string) error {
	if len(args) == 0 {
		honeyHelp()
		return nil
	}

	if args[0] == "help" {
		honeyHelp()
		return nil
	}

	if args[0] == "listfiletypes" {
		honeyListFileTypes()
		return nil
	}

	confPath := flag.String("conf", "", "path to config file")
	path := flag.String("path", "", "path to decoy folder")
	breadth := flag.Int("b", 1, "how many files roughly per folder")
	depth := flag.Int("d", 1, "roughly how deep to go from the root path provided")
	minsize := flag.Int("minsize", 1, "minimum size of each file in KB")
	maxsize := flag.Int("maxsize", 9, "maximum size of each file in KB")
	filetypes := flag.String("filetypes", "", "comma separated list of filetypes to generate")

	flag.CommandLine.Parse(args)

	conf := &honeyConf{}

	if *confPath == "" {
		if *path == "" {
			log.Fatalln("pattern option required")
		}
		conf.Path = *path
		conf.Breadth = *breadth
		conf.Depth = *depth
		conf.Filetypes = strings.Split(*filetypes, ",")
		conf.MaxSize = *maxsize
		conf.MinSize = *minsize
	} else {
		c, err := loadHoneyConf(*confPath)
		if err != nil {
			log.Fatalln("couldn't load conf:", *confPath)
		}
		conf = c
	}

	fmt.Println("P:", conf.Path)
	fmt.Println("D:", conf.Depth)
	fmt.Println("B:", conf.Breadth)
	fmt.Println("Max Size:", conf.MaxSize)
	fmt.Println("Min Size:", conf.MinSize)
	fmt.Println("Filetypes:", conf.Filetypes)

	return nil
}

func loadHoneyConf(confPath string) (*honeyConf, error) {
	data, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}

	var conf *honeyConf
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
