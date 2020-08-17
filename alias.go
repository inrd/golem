package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
)

var Aliases = make(map[string]string)

func RestoreAliases(aliasesFilePath string) {
	_, err := os.Stat(aliasesFilePath)
	if os.IsNotExist(err) {
		intro := "# Golem Aliases\n# Add each alias on a separate line :\n# [alias] [commands]\n# [alias] must be a single word but [commands] can have an arbitrary length\n"
		err = ioutil.WriteFile(aliasesFilePath, []byte(intro), 0644)
		if err != nil {
			PrintError(err.Error())
		}
	} else {
		f, err := os.Open(aliasesFilePath)
		if err != nil {
			PrintError(err.Error())
			return
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			if fields[0] != "#" {
				Aliases[fields[0]] = strings.Join(fields[1:], " ")
			}
		}

		err = scanner.Err()
		if err != nil {
			PrintError(err.Error())
		}
	}
}

func SubstituteAlias(args []string) []string {
	alias, exists := Aliases[args[0]]
	if exists {
		expr := strings.Fields(alias)
		return append(expr, args[1:]...)
	} else {
		return args
	}
}
