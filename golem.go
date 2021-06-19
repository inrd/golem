package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const version = "2.0.1"

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

	RestoreAliases(AliasesFilePath(golemDirPath))

	const prompt = "\n$> "
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf(prompt)
	for scanner.Scan() {
		ln := scanner.Text()
		if len(ln) > 0 {
			cmd := SubstituteAlias(strings.Fields(ln))
			running := true
			found := false

			if cmd[0] == "help" {
				HelpCommand(cmd[1:])
				found = true
			} else {
				for _, c := range Commands {
					if c.name == cmd[0] {

						if c.mutator {
							if c.name != "undo" {
								PushToHistory(tree.serialize())
							}
						}

						if c.name == "undo" {
							data := Undo()
							if data == nil {
								PrintError("Nothing to undo")
							} else {
								tmp := tree.deserialize(data)
								if tmp == nil {
									PrintError("Unable to restore state")
								} else {
									tree = tmp
									listFocusedNode(tree)
								}
							}
						} else {
							running = c.callback(cmd[1:], tree)
						}
						if c.mutator {
							tree.SaveToFile(golemTreeFilePath)
						}
						found = true
						break
					}
				}
			}
			if !found {
				PrintError("Invalid Command")
			}
			if !running {
				break
			}
		} else {
			PrintError("Enter a command (type 'help' for a list of commands)")
		}

		fmt.Printf("\n%v", prompt)

	}

	err = scanner.Err()
	if err != nil {
		log.Println(err)
	}
}
