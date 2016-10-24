package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/clagraff/argparse"
)

func createDB(path string) error {
	if out, err := exec.Command("sqlite3", path, ".databases").Output(); err != nil {
		return err
	} else if len(out) == 0 {
		return errors.New("Could not create database.")
	}
	return nil
}

func runInit(p *argparse.Parser, ns *argparse.Namespace, args []string, err error) {
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

	path := ns.Get("path")
	if len(path.(string)) == 0 {
		errorOut("p, path: must be a valid file name")
		return
	}

	fileExists := true

	stat, err := os.Stat(path.(string))
	if err != nil {
		if os.IsNotExist(err) == true {
			fileExists = false
		} else {
			errorOut(err.Error())
			return
		}
	}

	if fileExists == false {
		err = createDB(path.(string))
		if err != nil {
			errorOut("p, path: " + err.Error())
			return
		}
	} else {
		if stat.IsDir() == true {
			errorOut("p, path: a directory is not a database")
			return
		}
	}

	dir := ns.Get("dir")
	if len(dir.(string)) == 0 {
		errorOut("d, dir: must be a valid directory name")
		return
	}

	dirExists := true

	stat, err = os.Stat(dir.(string))
	if err != nil {
		if os.IsNotExist(err) == true {
			dirExists = false
		} else {
			errorOut(err.Error())
			return
		}
	}

	if dirExists == false {
		err = os.MkdirAll(dir.(string), 0777)
	}

	cfg := Config{}
	cfg.Path = path.(string)
	cfg.Dir = dir.(string)
	cfg.AddTrx = true

	err = WriteCfg(cfg)
	if err != nil {
		errorOut(err.Error())
	}
}

func AddInitParser(mainParser *argparse.Parser) {
	p := argparse.NewParser("ducky - init")
	p.AddHelp()

	dbName := argparse.NewArg("p Path", "path", "Path to sqlite database").Default("db.sqlite3").NotRequired()
	sqlDir := argparse.NewArg("d dir", "dir", "Path to migration directory").Default("sql").NotRequired()
	p.AddOptions(dbName, sqlDir)

	mainParser.AddParser("init", p, runInit)
}
