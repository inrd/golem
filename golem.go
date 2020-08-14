package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
    "strconv"
    "errors"
    "regexp"
)

const version = "2.0"


type Node struct {
    text string
    parent *Node
    children []*Node
}


func (n *Node) String() string {
    return n.text;
}


type Tree struct {
    root *Node
    focus *Node
}

func createNode (text string, parent *Node) *Node {
    n := Node{text: text, parent: parent}
    return &n
}

func initTree (txt string) *Tree {
    n := createNode(txt, nil)
    t := Tree{root: n, focus: n}
    return &t
}


func recurse (n *Node, callback func(*Node)) {
    if n.children != nil {
        for i,_ := range n.children {
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
        p = strings.Join([]string{"/", np.text, p}, "")
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

func parseIndexArgs (args []string, t *Tree) (int, error) {
    if args == nil {
        return -1, errors.New("No arguments")
    }
    idx, err := strconv.Atoi(args[0])
    if err != nil {
        // argument is a search pattern
        patt := strings.Join(args, " ")
        ipatt := strings.Join([]string{"(?i)", patt}, "")
        for i,c := range t.focus.children {
            match, err := regexp.MatchString(ipatt, c.text)
            if  err != nil {
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




func listFocusedNode (t *Tree) {
    fmt.Printf("\n%v\n", t.getFocusPath())
    for i,n := range t.focus.children {
        fmt.Println(i + 1, ". ", n)
    }
}



func addCommand (args []string, t *Tree) bool {
    txt := strings.Join(args, " ")
    t.focus.children = append(t.focus.children, createNode(txt, t.focus))
    listFocusedNode(t)
    return true
}


type Command struct {
	name     string
	help     string
	callback func([]string, *Tree) bool
}

var commands = []Command{
	{
		name:     "add",
        help:     "adds a node", //TODO: fix help text
		callback: addCommand,
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
		for i, c := range commands {
			if c.name == cmd[0] {
				running = commands[i].callback(cmd[1:], tree)
				break
			}
		}
		if !running {
			log.Println("Exiting...")
			break
		}
		fmt.Printf("\n%v", prompt)
	}

	err := scanner.Err()
	if err != nil {
		log.Println(err)
	}
}
