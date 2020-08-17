package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
		name:     "tree",
		summary:  "append new node and focus it",
		help:     ":::: tree [text] ::::\nappenda new node to the current list with the content defined in [text] and move focus to the new node",
		callback: treeCommand,
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
		name:     "move",
		summary:  "move a node inside the current list",
		help:     ":::: move [node index || search pattern] [target index] ::::\nmove the node at [node index] or the first node matching [serach pattern] to [target index] in the current list, shifting the other nodes",
		callback: moveCommand,
		mutator:  true,
	},
	{
		name:     "cut",
		summary:  "cut node at index",
		help:     ":::: cut [node index || search pattern] ::::\ndelete the node at [node index] or containing [search pattern] and copy its contents (recursively) to the internal clipboard",
		callback: cutCommand,
		mutator:  true,
	},
	{
		name:     "paste",
		summary:  "paste node",
		help:     ":::: paste ::::\nappend the contents of the internal clipboard to the current list\nPasting a node purges the clipboard.",
		callback: pasteCommand,
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
		name:     "search",
		summary:  "search expression in the whole tree",
		help:     ":::: search [regexp] ::::\nsearch (case insensitive) the whole tree for [regexp]",
		callback: searchCommand,
		mutator:  false,
	},
	{
		name:     "undo",
		summary:  "undo previous destructive action",
		help:     ":::: undo ::::\nrevert latest destructive action (will reset the focus to the root node)",
		callback: nil, // special case : handled in main()
		mutator:  true,
	},
	{
		name:     "export",
		summary:  "export plain text subtree of the current node",
		help:     "::::export [filename]::::\nwrite a plain text version of the currently focused node to [filename]",
		callback: exportCommand,
		mutator:  false,
	},
	{
		name:     "alias",
		summary:  "list registered command aliases",
		help:     ":::: alias ::::\nprint a list of the registered aliases\naliases can be edited in $HOME/.golem/aliases",
		callback: aliasCommand,
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
	fmt.Printf("\n%v", msg)
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

func treeCommand(args []string, t *Tree) bool {
	if checkArgs(args, 1) == false {
		return true
	}
	txt := strings.Join(args, " ")
	t.focus.Children = append(t.focus.Children, CreateNode(txt, t.focus))
	t.focus = t.focus.Children[len(t.focus.Children)-1]
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

func moveCommand(args []string, t *Tree) bool {
	if checkArgs(args, 2) == false {
		return true
	}
	posArg := args[:1]
	idx, perr := parseIndexArgs(posArg, t)
	dest, oerr := strconv.Atoi(args[1])

	if perr != nil {
		PrintError(perr.Error())
		return true
	}
	if oerr != nil {
		PrintError(oerr.Error())
		return true
	}
	if idx < 0 || idx >= len(t.focus.Children) {
		PrintError("Invalid node")
		return true
	}

	dest -= 1
	if dest < 0 {
		dest = 0
	} else if dest >= len(t.focus.Children) {
		dest = len(t.focus.Children) - 1
	}

	if idx < dest {
		tmp := append([]*Node{}, t.focus.Children[:idx]...)
		tmp = append(tmp, t.focus.Children[idx+1:dest+1]...)
		tmp = append(tmp, t.focus.Children[idx])
		tmp = append(tmp, t.focus.Children[dest+1:]...)
		t.focus.Children = tmp
	} else if idx > dest {
		tmp := append([]*Node{}, t.focus.Children[:dest]...)
		tmp = append(tmp, t.focus.Children[idx])
		tmp = append(tmp, t.focus.Children[dest:idx]...)
		tmp = append(tmp, t.focus.Children[idx+1:]...)
		t.focus.Children = tmp
	}
	listFocusedNode(t)
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

func recursiveSearch(n *Node, pattern *regexp.Regexp, results *[]string, callback func(*Node, *regexp.Regexp, *[]string)) {
	if n.Children != nil {
		for i, _ := range n.Children {
			recursiveSearch(n.Children[i], pattern, results, callback)
		}
	}
	callback(n, pattern, results)
}

func searchNode(n *Node, pattern *regexp.Regexp, results *[]string) {
	match := pattern.MatchString(n.Text)
	if match {
		var s string
		if n.parent == nil {
			s = fmt.Sprintf("/%v", n.Text)
		} else {
			t := n
			for t.parent != nil {
				s = fmt.Sprintf("/%v%v", t.Text, s)
				t = t.parent
			}
		}
		*results = append(*results, fmt.Sprintf("=> %v", s))
	}
}

func searchCommand(args []string, t *Tree) bool {
	if checkArgs(args, 1) == false {
		return true
	}
	var results []string
	pattern := strings.Join(args, " ")
	ipatt := fmt.Sprintf("(?i)%v", pattern)
	re, err := regexp.Compile(ipatt)
	if err != nil {
		PrintError(err.Error())
		return true
	}
	recursiveSearch(t.Root, re, &results, searchNode)
	if len(results) < 1 {
		PrintError("Pattern not found")
	} else {
		fmt.Printf("\n%v occurences found : \n\n", len(results))
		for _, s := range results {
			fmt.Println(s)
		}
	}

	return true
}

// BUG : does not work properly from root node (seems to work fine otherwise)
func recursiveDump(n *Node, txt *string, focus *Node, callback func(*Node, *string, *Node)) {
	if n.Children != nil {
		for i, _ := range n.Children {
			recursiveDump(n.Children[i], txt, focus, callback)
		}
	}
	callback(n, txt, focus)
}

func dumpNode(n *Node, txt *string, focus *Node) {
	if n != focus {
		depth := 0
		np := n
		p := np.parent
		for p != focus && p != nil {
			depth += 1
			np = p
			p = np.parent
		}
		for i := 0; i < depth; i++ {
			*txt = fmt.Sprintf("%v  ", *txt)
		}
		*txt = fmt.Sprintf("%v- %v\n", *txt, n.Text)
	}
}

func exportCommand(args []string, t *Tree) bool {

	if checkArgs(args, 1) == false {
		return true
	}

	var dump string
	recursiveDump(t.focus, &dump, t.focus, dumpNode)

	dir, err := os.Getwd()
	if err != nil {
		PrintError(err.Error())
		return true
	}

	path := filepath.Join(dir, args[0])

	err = ioutil.WriteFile(path, []byte(dump), 0644)
	if err != nil {
		PrintError(err.Error())
	} else {
		fmt.Printf("\nWrote plain text subtree to %v", path)
	}
	return true
}

func HelpCommand(args []string) {
	if len(args) > 0 {
		found := false
		for _, c := range Commands {
			if args[0] == c.name {
				fmt.Printf("\n%v", c.help)
				found = true
				break
			}
		}
		if found == false {
			PrintError("Invalid command name")
		}
	} else {
		fmt.Printf("\nType 'help [command]' to get help about a specific command.\n\nPossible commands :\n\n")
		for _, c := range Commands {
			fmt.Printf("* %v => %v\n", c.name, c.summary)
		}
	}
}

var clipboard *Node

func cutCommand(args []string, t *Tree) bool {
	if checkArgs(args, 1) == false {
		return true
	}
	idx, err := parseIndexArgs(args, t)
	if err != nil {
		PrintError(err.Error())
		return true
	}
	if idx < 0 || idx >= len(t.focus.Children) {
		PrintError("Invalid node")
		return true
	} else {
		clipboard = t.focus.Children[idx]
		clipboard.parent = nil
		t.focus.Children = append(t.focus.Children[:idx], t.focus.Children[idx+1:]...)
		fmt.Printf("\nNode [%v] copied to the internal clipboard\n", clipboard.Text)
		listFocusedNode(t)
	}
	return true
}

func pasteCommand(args []string, t *Tree) bool {
	if clipboard == nil {
		PrintError("Nothing to paste...")
		return true
	}
	clipboard.parent = t.focus
	t.focus.Children = append(t.focus.Children, clipboard)
	clipboard = nil
	listFocusedNode(t)
	return true
}

func aliasCommand([]string, *Tree) bool {
	if len(Aliases) == 0 {
		PrintError("No aliases registered")
	} else {
		fmt.Printf("\n")
		for al, cmd := range Aliases {
			fmt.Printf("* %v => %v\n", al, cmd)
		}
	}
	return true
}
