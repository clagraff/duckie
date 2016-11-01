package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/clagraff/argparse"
)

func runStatus(p *argparse.Parser, ns *argparse.Namespace, args []string, err error) {
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

	var completed []string
	var pending []string

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

	counter := make(map[int]int)
	for _, f := range files {
		splitted := strings.Split(filepath.Base(f), "_")
		if fID, err := strconv.Atoi(splitted[0]); err == nil {
			if fID > version {
				counter[fID]++
				if counter[fID] > 1 {
					errorOut(fmt.Sprintf("Multiple migrations share version number: %d, file: %s", fID, f))
					return
				}
				pending = append(pending, f)
			} else {
				completed = append(completed, f)
			}
		}
	}

	fmt.Printf(
		"Current version: %d, Completed Migrations: %d, Pending Migrations: %d\n",
		version,
		len(completed),
		len(pending),
	)
	if len(pending)+len(completed) == 0 {
		fmt.Println("No migrations")
	} else {
		fmt.Printf("\n")
	}
	for _, f := range completed {
		fmt.Println("[DONE]   ", f)
	}
	for _, f := range pending {
		fmt.Println("[PENDING]", f)
	}
}

func AddStatusParser(mainParser *argparse.Parser) {
	p := argparse.NewParser("ducky - status", runStatus)
	p.AddHelp()

	mainParser.AddParser("status", p)
}
