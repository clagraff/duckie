package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/clagraff/argparse"
)

func checkSqlite() error {
	if out, err := exec.Command("which", "sqlite3").Output(); err != nil {
		return err
	} else if len(out) == 0 {
		return errors.New("Could not locate `sqlite3` command.")
	}
	return nil
}

func main() {
	p := argparse.NewParser("ducky").Version("0.0.0")
	p.AddHelp().AddVersion() // Enable help and version flags

	if checkSqlite() != nil {
		fmt.Println("Could not locate `sqlite3` command. Please install `sqlite3`.\n")
		p.ShowHelp()
		return
	}

	AddInitParser(p)

	// Parse all available program arguments (except for the program path).
	if _, _, err := p.Parse(os.Args[1:]...); err != nil {
		switch err.(type) {
		case argparse.ShowHelpErr, argparse.ShowVersionErr:
			return
		default:
			fmt.Println(err, "\n")
			p.ShowHelp()
		}
	}
}
