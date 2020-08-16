package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var Commands = []Command{
	{
		name: "list",
		help: "prints contents of current node",
		callback: func(args []string, t *Tree) bool {
			listFocusedNode(t)
			return true
		},
		mutator: false,
	},
	{
		name:     "add",
		help:     "adds a node", //TODO: fix help text
		callback: addCommand,
		mutator:  true,
	},
	{
		name:     "del",
		help:     "deletes a node", //TODO: fix help text
		callback: delCommand,
		mutator:  true,
	},
	{
		name:     "focus",
		help:     "focus a node", //TODO: fix help text
		callback: focusCommand,
		mutator:  false,
	},
	{
		name:     "back",
		help:     "go back in tree", //TODO: fix help text
		callback: backCommand,
		mutator:  false,
	},
	{
		name:     "root",
		help:     "go to tree root", //TODO: fix help text
		callback: rootCommand,
		mutator:  false,
	},
	{
		name: "exit",
		help: "quits the program",
		callback: func([]string, *Tree) bool {
			return false
		},
		mutator: false,
	},
}

func PrintError(msg string) {
	fmt.Printf("\n%v\n", msg)
}

func parseIndexArgs(args []string, t *Tree) (int, error) {
	if args == nil {
		return -1, errors.New("No arguments")
	}
	idx, err := strconv.Atoi(args[0])
	if err != nil {
		// argument is a search pattern
		patt := strings.Join(args, " ")
		ipatt := fmt.Sprintf("(?i)%v", patt)
		for i, c := range t.Focus().Children {
			match, err := regexp.MatchString(ipatt, c.Text)
			if err != nil {
				return -1, err
			}
			if match {
				return i, nil
			}
		}
		return -1, errors.New("Pattern not found")
	} else {
		return idx - 1, nil
	}
}

func listFocusedNode(t *Tree) {
	fmt.Printf("\n::::%v\n", t.GetFocusPath())
	if len(t.Focus().Children) == 0 {
		PrintError("Empty node...")
	} else {
		for i, n := range t.Focus().Children {
			fmt.Println(i+1, ". ", n)
		}
	}
}

func addCommand(args []string, t *Tree) bool {
	txt := strings.Join(args, " ")
	t.Focus().Children = append(t.Focus().Children, CreateNode(txt, t.Focus()))
	listFocusedNode(t)
	return true
}

func delCommand(args []string, t *Tree) bool {
	idx, err := parseIndexArgs(args, t)
	if err != nil {
		PrintError(err.Error())
	} else {
		if idx >= 0 && idx < len(t.focus.Children) {
			t.focus.Children = append(t.focus.Children[:idx], t.focus.Children[idx+1:]...)
			listFocusedNode(t)
		} else {
			PrintError("Inavlid node")
		}
	}
	return true
}

func focusCommand(args []string, t *Tree) bool {
	idx, err := parseIndexArgs(args, t)
	if err != nil {
		PrintError(err.Error())
	} else {
		if idx >= 0 && idx < len(t.focus.Children) {
			t.focus = t.focus.Children[idx]
			listFocusedNode(t)
		} else {
			PrintError("Inavlid node")
		}
	}
	return true
}

func backCommand(args []string, t *Tree) bool {
	if t.focus.parent != nil {
		t.focus = t.focus.parent
		listFocusedNode(t)
	} else {
		PrintError("Already at root...")
	}
	return true
}

func rootCommand(args []string, t *Tree) bool {
	if t.focus.parent != nil {
		t.focus = t.Root
		listFocusedNode(t)
	} else {
		PrintError("Already at root...")
	}
	return true
}

type Command struct {
	name     string
	help     string
	callback func([]string, *Tree) bool
	mutator  bool
}
