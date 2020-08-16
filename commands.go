package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Command struct {
	name     string
	summary  string
	help     string
	callback func([]string, *Tree) bool
	mutator  bool
}

var Commands = []Command{
	{
		name:    "list",
		summary: "list contents of current node",
		help:    ":::: list ::::\nprint the contents of the current node to the console\nnodes which contain sub-nodes have a little indicator ([+]) after their index",
		callback: func(args []string, t *Tree) bool {
			listFocusedNode(t)
			return true
		},
		mutator: false,
	},
	{
		name:     "add",
		summary:  "append new node",
		help:     ":::: add [text] ::::\nappend a new node to the current list with the content defined in [text]",
		callback: addCommand,
		mutator:  true,
	},
	{
		name:     "insert",
		summary:  "insert new node at index",
		help:     ":::: insert [index] [text] ::::\ninsert a new node in the current list at [index] with the content defined in [text]",
		callback: insertCommand,
		mutator:  true,
	},
	{
		name:     "edit",
		summary:  "edit node at index",
		help:     ":::: edit [node index] [new text] ::::\nreplace the node text at [node index] by [new text]",
		callback: editCommand,
		mutator:  true,
	},
	{
		name:     "del",
		summary:  "delete node at index",
		help:     ":::: del [node index] ::::\ndelete the node at [node index]",
		callback: delCommand,
		mutator:  true,
	},
	{
		name:     "focus",
		summary:  "bring focus to a specific node",
		help:     ":::: focus [node index || search pattern] ::::\nset the root node to [node index] or to the first node matching [search pattern]",
		callback: focusCommand,
		mutator:  false,
	},
	{
		name:     "back",
		summary:  "move up in the tree",
		help:     ":::: back ::::\nmove back (up) in the hierarchy",
		callback: backCommand,
		mutator:  false,
	},
	{
		name:     "root",
		summary:  "move to the tree root",
		help:     ":::: root ::::\nmove back to the root node",
		callback: rootCommand,
		mutator:  false,
	},
	{
		name:    "exit",
		summary: "exit Golem",
		help:    ":::: exit ::::\nexit Golem",
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
		for i, c := range t.focus.Children {
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

func checkArgs(args []string, minLength int) bool {
	if len(args) < minLength {
		PrintError("Missing arguments")
		return false
	} else {
		return true
	}
}

func listFocusedNode(t *Tree) {
	fmt.Printf("\n::::%v\n", t.GetFocusPath())
	if len(t.focus.Children) == 0 {
		PrintError("Empty node...")
	} else {
		for i, n := range t.focus.Children {
			fmt.Println(i+1, ". ", n)
		}
	}
}

func addCommand(args []string, t *Tree) bool {
	if checkArgs(args, 1) == false {
		return true
	}
	txt := strings.Join(args, " ")
	t.focus.Children = append(t.focus.Children, CreateNode(txt, t.focus))
	listFocusedNode(t)
	return true
}

func insertCommand(args []string, t *Tree) bool {
	if checkArgs(args, 2) == false {
		return true
	}
	posArg := args[:1]
	txt := strings.Join(args[1:], " ")
	idx, err := parseIndexArgs(posArg, t)
	if err != nil {
		PrintError(err.Error())
	} else {
		if idx <= 0 {
			tmp := t.focus.Children
			t.focus.Children = append([]*Node{}, CreateNode(txt, t.focus))
			t.focus.Children = append(t.focus.Children, tmp...)
		} else if idx >= len(t.focus.Children) {
			t.focus.Children = append(t.focus.Children, CreateNode(txt, t.focus))
		} else {
			tmp := append([]*Node{}, t.focus.Children[:idx]...)
			tmp = append(tmp, CreateNode(txt, t.focus))
			tmp = append(tmp, t.focus.Children[idx:]...)
			t.focus.Children = tmp
		}
		listFocusedNode(t)
	}
	return true
}

func editCommand(args []string, t *Tree) bool {
	if checkArgs(args, 2) == false {
		return true
	}
	posArg := args[:1]
	txt := strings.Join(args[1:], " ")
	idx, err := parseIndexArgs(posArg, t)
	if err != nil {
		PrintError(err.Error())
	} else {
		if idx >= 0 && idx < len(t.focus.Children) {
			t.focus.Children[idx].Text = txt
			listFocusedNode(t)
		} else {
			PrintError("Inavlid node")
		}
	}
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
	if checkArgs(args, 1) == false {
		return true
	}
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
