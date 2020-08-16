package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const version = "2.0"

/*TODO
type Alias struct {
    name string
    target string
}
*/

func main() {
	var tree *Tree

	fmt.Println("\n::::Golem v", version)

	golemDirPath := InitGolemDir()
	golemTreeFilePath := TreeFilePath(golemDirPath)

	_, err := os.Stat(golemTreeFilePath)
	if os.IsNotExist(err) {
		tree = InitTree("Golem")
		tree.SaveToFile(golemTreeFilePath)
	} else {
		tree = tree.RestoreFromFile(golemTreeFilePath)
		if tree == nil {
			tree = InitTree("Golem")
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

		for _, c := range Commands {
			if c.name == cmd[0] {
				running = c.callback(cmd[1:], tree)
				if c.mutator {
					tree.SaveToFile(golemTreeFilePath)
				}
				found = true
				break
			}
		}
		if !found {
			PrintError("Invalid Command")
		}
		if !running {
			fmt.Println("\nGoodbye!\n")
			break
		}
		fmt.Printf("\n%v", prompt)
	}

	err = scanner.Err()
	if err != nil {
		log.Println(err)
	}
}
