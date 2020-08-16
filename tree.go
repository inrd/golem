package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

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

func CreateNode(text string, parent *Node) *Node {
	n := Node{Text: text, parent: parent}
	return &n
}

func InitTree(txt string) *Tree {
	n := CreateNode(txt, nil)
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

func (t *Tree) Traverse(callback func(*Node)) {
	recurse(t.Root, callback)
}

func (t *Tree) initParents() {
	recurseAndSetParent(t.Root, nil)
}

func (t *Tree) serialize() []byte {
	data, err := json.Marshal(t)
	if err != nil {
		PrintError(err.Error())
		return nil
	}
	return data
}

func (t *Tree) deserialize(data []byte) *Tree {
	tmp := Tree{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		PrintError(err.Error())
		return nil
	}
	tmp.initParents()
	tmp.focus = tmp.Root
	return &tmp
}

func (t *Tree) GetFocusPath() string {
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

func (t *Tree) SaveToFile(path string) {
	data := t.serialize()
	if data != nil {
		err := ioutil.WriteFile(path, t.serialize(), 0644)
		if err != nil {
			PrintError(err.Error())
		}
	}
}

func (t *Tree) RestoreFromFile(path string) *Tree {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		PrintError(err.Error())
		return nil
	} else {
		tp := t.deserialize(data)
		if tp == nil {
			PrintError("Error parsing data")
			return nil
		}
		return tp
	}
}
