package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const version = "2.0"

type Node struct {
	text     string
	parent   *Node
	children []*Node
}

func (n *Node) String() string {
	if len(n.children) > 0 {
		return fmt.Sprintf("[+] %v", n.text)
	}
	return n.text
}

type Tree struct {
	root  *Node
	focus *Node
}

func createNode(text string, parent *Node) *Node {
	n := Node{text: text, parent: parent}
	return &n
}

func initTree(txt string) *Tree {
	n := createNode(txt, nil)
	t := Tree{root: n, focus: n}
	return &t
}

func recurse(n *Node, callback func(*Node)) {
	if n.children != nil {
		for i, _ := range n.children {
			recurse(n.children[i], callback)
		}
	}
	callback(n)
}

func (t *Tree) traverse(callback func(*Node)) {
	recurse(t.root, callback)
}

func (t *Tree) getFocusPath() string {
	if t.focus.parent == nil {
		return "/"
	}
	p := ""
	np := t.focus
	for np.parent != nil {
		p = fmt.Sprintf("/%v%v", np.text, p)
		np = np.parent
	}
	return p
}

/*TODO
type Alias struct {
    name string
    target string
}
*/

func printError(msg string) {
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
		for i, c := range t.focus.children {
			match, err := regexp.MatchString(ipatt, c.text)
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
	fmt.Printf("\n::::%v\n", t.getFocusPath())
	if len(t.focus.children) == 0 {
		printError("Empty node...")
	} else {
		for i, n := range t.focus.children {
			fmt.Println(i+1, ". ", n)
		}
	}
}

func addCommand(args []string, t *Tree) bool {
	txt := strings.Join(args, " ")
	t.focus.children = append(t.focus.children, createNode(txt, t.focus))
	listFocusedNode(t)
	return true
}

func focusCommand(args []string, t *Tree) bool {
	idx, err := parseIndexArgs(args, t)
	if err != nil {
		printError(err.Error())
	} else {
		if idx >= 0 && idx < len(t.focus.children) {
			t.focus = t.focus.children[idx]
			listFocusedNode(t)
		} else {
			printError("Inavlid node")
		}
	}
	return true
}

func backCommand(args []string, t *Tree) bool {
	if t.focus.parent != nil {
		t.focus = t.focus.parent
		listFocusedNode(t)
	} else {
		printError("Already at root...")
	}
	return true
}

func rootCommand(args []string, t *Tree) bool {
	if t.focus.parent != nil {
		t.focus = t.root
		listFocusedNode(t)
	} else {
		printError("Already at root...")
	}
	return true
}

type Command struct {
	name     string
	help     string
	callback func([]string, *Tree) bool
}

var commands = []Command{
	{
		name: "list",
		help: "prints contents of current node",
		callback: func(args []string, t *Tree) bool {
			listFocusedNode(t)
			return true
		},
	},
	{
		name:     "add",
		help:     "adds a node", //TODO: fix help text
		callback: addCommand,
	},
	{
		name:     "focus",
		help:     "focus a node", //TODO: fix help text
		callback: focusCommand,
	},
	{
		name:     "back",
		help:     "go back in tree", //TODO: fix help text
		callback: backCommand,
	},
	{
		name:     "root",
		help:     "go to tree root", //TODO: fix help text
		callback: rootCommand,
	},
	{
		name: "exit",
		help: "quits the program",
		callback: func([]string, *Tree) bool {
			return false
		},
	},
}

func main() {
	tree := initTree("Golem")
	//TODO: var aliases []Alias

	const prompt = "$> "
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf(prompt)
	for scanner.Scan() {
		ln := scanner.Text()
		cmd := strings.Fields(ln)
		running := true
		found := false

		for i, c := range commands {
			if c.name == cmd[0] {
				running = commands[i].callback(cmd[1:], tree)
				found = true
				break
			}
		}
		if !found {
			printError("Invalid Command")
		}
		if !running {
			fmt.Println("Goodbye!")
			break
		}
		fmt.Printf("\n%v", prompt)
	}

	err := scanner.Err()
	if err != nil {
		log.Println(err)
	}
}
