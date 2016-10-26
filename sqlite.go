package main

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func MigrationText(cfg Config) string {
	var text []string

	text = append(text, "-- @ducky Up")
	if cfg.AddTxn == true {
		text = append(
			text,
			"BEGIN DEFERRED TRANSACTION;",
			"\n",
			"COMMIT TRANSACTION;",
		)
	}

	text = append(text, "\n\n\n-- @ducky Down")
	if cfg.AddTxn == true {
		text = append(
			text,
			"BEGIN DEFERRED TRANSACTION;",
			"\n",
			"COMMIT TRANSACTION;",
		)
	}
	text = append(text, "\n")

	return strings.Join(text, "\n")
}

func GetUserVersion(cfg Config) (int, error) {
	if _, err := os.Stat(cfg.Path); err != nil {
		if os.IsNotExist(err) {
			return 0, errors.New("Could not find database")
		}
	}

	if out, err := exec.Command("sqlite3", cfg.Path, "PRAGMA user_version;").Output(); err != nil {
		return 0, err
	} else if len(out) == 0 {
		return 0, errors.New("Could not determine database user version")
	} else {
		id, err := strconv.Atoi(string(out[:len(out)-1]))
		if err != nil {
			return 0, errors.New("sqlite error: " + string(out))
		}
		return id, nil
	}
}

func WriteUserVersion(version int, cfg Config) error {
	if _, err := os.Stat(cfg.Path); err != nil {
		if os.IsNotExist(err) {
			return errors.New("Could not find database")
		}
	}

	if out, err := exec.Command("sqlite3", cfg.Path, "PRAGMA user_version="+strconv.Itoa(version)+";").CombinedOutput(); err != nil {
		return err
	} else if len(out) != 0 {
		return errors.New("sqlite error: " + string(out))
	}

	if id, err := GetUserVersion(cfg); err != nil {
		return err
	} else if id != version {
		return errors.New("Failed to set database user version")
	}

	return nil
}
