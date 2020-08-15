package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const version = "2.0"

type Node struct {
	Text     string
	parent   *Node
	Children []*Node
}

func (n *Node) String() string {
	if len(n.Children) > 0 {
		return fmt.Sprintf("[+] %v", n.Text)
	}
	return n.Text
}

type Tree struct {
	Root  *Node
	focus *Node
}

func createNode(text string, parent *Node) *Node {
	n := Node{Text: text, parent: parent}
	return &n
}

func initTree(txt string) *Tree {
	n := createNode(txt, nil)
	t := Tree{Root: n, focus: n}
	return &t
}

func recurse(n *Node, callback func(*Node)) {
	if n.Children != nil {
		for i, _ := range n.Children {
			recurse(n.Children[i], callback)
		}
	}
	callback(n)
}

func recurseAndSetParent(n *Node, p *Node) {
	if n.Children != nil {
		for i, _ := range n.Children {
			recurseAndSetParent(n.Children[i], n)
		}
	}
	n.parent = p
}

func (t *Tree) traverse(callback func(*Node)) {
	recurse(t.Root, callback)
}

func (t *Tree) initParents() {
	recurseAndSetParent(t.Root, nil)
}

func (t *Tree) serialize() []byte {
	data, err := json.Marshal(t)
	if err != nil {
		printError(err.Error())
		return nil
	}
	return data
}

func (t *Tree) deserialize(data []byte) *Tree {
	tmp := Tree{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		printError(err.Error())
		return nil
	}
	tmp.initParents()
	tmp.focus = tmp.Root
	return &tmp
}

func (t *Tree) getFocusPath() string {
	if t.focus.parent == nil {
		return "/"
	}
	p := ""
	np := t.focus
	for np.parent != nil {
		p = fmt.Sprintf("/%v%v", np.Text, p)
		np = np.parent
	}
	return p
}

func (t *Tree) saveToFile(path string) {
	data := t.serialize()
	if data != nil {
		err := ioutil.WriteFile(path, t.serialize(), 0644)
		if err != nil {
			printError(err.Error())
		}
	}
}

func (t *Tree) restoreFromFile(path string) *Tree {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		printError(err.Error())
		return nil
	} else {
		tp := t.deserialize(data)
		if tp == nil {
			printError("Error parsing data")
			return nil
		}
		return tp
	}
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

func listFocusedNode(t *Tree) {
	fmt.Printf("\n::::%v\n", t.getFocusPath())
	if len(t.focus.Children) == 0 {
		printError("Empty node...")
	} else {
		for i, n := range t.focus.Children {
			fmt.Println(i+1, ". ", n)
		}
	}
}

func addCommand(args []string, t *Tree) bool {
	txt := strings.Join(args, " ")
	t.focus.Children = append(t.focus.Children, createNode(txt, t.focus))
	listFocusedNode(t)
	return true
}

func delCommand(args []string, t *Tree) bool {
	idx, err := parseIndexArgs(args, t)
	if err != nil {
		printError(err.Error())
	} else {
		if idx >= 0 && idx < len(t.focus.Children) {
            t.focus.Children = append(t.focus.Children[:idx], t.focus.Children[idx + 1:]...)
			listFocusedNode(t)
		} else {
			printError("Inavlid node")
		}
	}
	return true
}

func focusCommand(args []string, t *Tree) bool {
	idx, err := parseIndexArgs(args, t)
	if err != nil {
		printError(err.Error())
	} else {
		if idx >= 0 && idx < len(t.focus.Children) {
			t.focus = t.focus.Children[idx]
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
		t.focus = t.Root
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
	mutator  bool
}

var commands = []Command{
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
        name: "del",
        help: "deletes a node", //TODO: fix help text
        callback: delCommand,
        mutator: true,
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

func main() {
	var tree *Tree

	const golemDirName = ".golem"
	golemDirPath := filepath.Join(os.Getenv("HOME"), golemDirName)
	_, derr := os.Stat(golemDirPath)
	if os.IsNotExist(derr) {
		dcerr := os.Mkdir(golemDirPath, 0755)
		if dcerr != nil {
			printError(dcerr.Error())
		}
	}

	const golemTreeFileName = "tree.json"
	golemTreeFilePath := filepath.Join(golemDirPath, golemTreeFileName)

	fmt.Println("\n::::Golem v", version)

	_, perr := os.Stat(golemTreeFilePath)
	if os.IsNotExist(perr) {
		tree = initTree("Golem")
		tree.saveToFile(golemTreeFilePath)
	} else {
		tree = tree.restoreFromFile(golemTreeFilePath)
		if tree == nil {
			tree = initTree("Golem")
		} else {
			listFocusedNode(tree)
		}
	}

	//TODO: var aliases []Alias

	const prompt = "\n$> "
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
				if c.mutator {
					tree.saveToFile(golemTreeFilePath)
				}
				found = true
				break
			}
		}
		if !found {
			printError("Invalid Command")
		}
		if !running {
			fmt.Println("\nGoodbye!\n")
			break
		}
		fmt.Printf("\n%v", prompt)
	}

	err := scanner.Err()
	if err != nil {
		log.Println(err)
	}
}
