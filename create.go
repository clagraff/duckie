package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/clagraff/argparse"
)

func runCreate(p *argparse.Parser, ns *argparse.Namespace, args []string, err error) {
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

	name := ns.Get("name").(string)

	cfg, err := ReadCfg()
	if err != nil {
		errorOut(err.Error())
		return
	}

	text := MigrationText(cfg)
	version, err := GetUserVersion(cfg)
	if err != nil {
		errorOut(err.Error())
		return
	}

	strVersion := fmt.Sprintf("%03d", version+1)
	err = ioutil.WriteFile(filepath.Join(cfg.Dir, strVersion+"_"+name+".sql"), []byte(text), 0777)
	if err != nil {
		errorOut(err.Error())
		return
	}
}

func AddCreateParser(mainParser *argparse.Parser) {
	p := argparse.NewParser("ducky - create", runCreate)
	p.AddHelp()

	name := argparse.NewArg("n name", "name", "Name of migration").Required()
	p.AddOption(name)

	mainParser.AddParser("create", p)
}
