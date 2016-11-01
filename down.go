package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/clagraff/argparse"
)

func runDown(p *argparse.Parser, ns *argparse.Namespace, args []string, err error) {
	if err != nil {
		switch err.(type) {
		case argparse.ShowHelpErr, argparse.ShowVersionErr:
			return
		default:
			fmt.Println(err, "\n")
			p.ShowHelp()
			return
		}
	}
	errorOut := func(msg string) {
		fmt.Println(msg, "\n")
		p.ShowHelp()
	}

	var all bool
	if ns.Get("all") == "true" {
		all = true
	}

	var filePaths []string

	cfg, err := ReadCfg()
	if err != nil {
		errorOut(err.Error())
		return
	}

	version, err := GetUserVersion(cfg)
	if err != nil {
		errorOut(err.Error())
		return
	}

	files, err := filepath.Glob(filepath.Join(cfg.Dir, "[0-9][0-9][0-9]_*.sql"))
	if err != nil {
		panic(err)
	}

	sortedFiles := sort.Reverse(sort.StringSlice(files))
	sort.Sort(sortedFiles)

	if all == true {
		files = []string{files[0]}
	}

	counter := make(map[int]int)
	for _, f := range files {
		splitted := strings.Split(filepath.Base(f), "_")
		if fID, err := strconv.Atoi(splitted[0]); err == nil {
			if fID <= version {
				counter[fID]++
				if counter[fID] > 1 {
					fmt.Println("Multiple migrations for same version:", strconv.Itoa(fID))
				}
				filePaths = append(filePaths, f)
			}
		}
	}

	fmt.Printf("Current version: %d, Target version: %d, Num of Migrations: %d\n", version, version-len(filePaths), len(filePaths))
	if len(filePaths) == 0 {
		fmt.Println("No migrations to run")
	}
	for _, path := range filePaths {
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		re := regexp.MustCompile(`-- @ducky Up((?:.|\n)*)-- @ducky Down((?:.|\n)*)`)
		matches := re.FindStringSubmatch(string(contents))
		if len(matches) != 3 {
			panic("bad migration file")
		}

		downMig := strings.TrimSpace(matches[2])

		if out, err := exec.Command("sqlite3", cfg.Path, downMig).CombinedOutput(); err != nil {
			fmt.Println("[FAILED]", path)
			fmt.Println("\n", string(out))
			return
		} else {
			fmt.Println("[OK]", path)
			version--
			err = WriteUserVersion(version, cfg)
			if err != nil {
				panic(err)
			}
		}
	}

}

func AddDownParser(mainParser *argparse.Parser) {
	p := argparse.NewParser("ducky - down", runDown)
	p.AddHelp()

	all := argparse.NewFlag("a all", "all", "Perform all possible downgrades from the current database version")
	p.AddOption(all)

	mainParser.AddParser("down", p)
}
