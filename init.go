package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
		fmt.Printf("Created database: %s\n", path.(string))
	} else {
		if stat.IsDir() == true {
			errorOut("p, path: a directory is not a database")
			return
		}
		fmt.Printf("Using existing database: %s\n", path.(string))
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
		if err != nil {
			errorOut(err.Error())
			return
		}
		fmt.Printf("Created migration directory: %s\n", dir.(string))
	} else {
		fmt.Printf("Using migration directory: %s\n", dir.(string))
	}

	cfg := Config{}
	cfg.Path = path.(string)
	cfg.Dir = dir.(string)

	txnStr := ns.Get("txn")
	if strings.ToLower(txnStr.(string)) == "true" {
		cfg.AddTxn = true
	} else {
		cfg.AddTxn = false
	}

	err = WriteCfg(cfg)
	if err != nil {
		errorOut(err.Error())
	}

	if err = WriteUserVersion(0, cfg); err != nil {
		errorOut(err.Error())
	}
}

func AddInitParser(mainParser *argparse.Parser) {
	p := argparse.NewParser("ducky - init", runInit)
	p.AddHelp()

	dbName := argparse.NewArg("p Path", "path", "Path to sqlite database").Default("database.sqlite3").NotRequired()
	sqlDir := argparse.NewArg("d dir", "dir", "Path to migration directory").Default("sql").NotRequired()
	addTxn := argparse.NewOption("t txn", "txn", "Do not auto-add transactions when creating migrations").Default("false")
	p.AddOptions(addTxn, dbName, sqlDir)

	mainParser.AddParser("init", p)
}
