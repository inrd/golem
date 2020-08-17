# :::: Golem  

A command driven outliner.

## Installation

If you have installed Go  : `go build && go install`

## Basic Usage

Run `golem` from the command line. It will open a prompt where you can enter commands to manipulate the tree.

Inside Golem, run `help` to get an overview of the available commands.

Some basic commands to get started :

- `add some item` : add a new node with the text "some item" to the current list.
- `focus 1` : move the focus to the node at index 1. Displays the subtree of node 1.
- `back`: move up by one level in the tree hierarchy.
- `root`: go back to the root node of the tree (where the program starts by default).
- `del 3` : delete the node at index 3.
- `list`: list content of the currently focused node
- `search something`: search the expression "something" in the whole tree.
- `undo`: undo previous destructive action.
- `export tree.txt` : export the focused tree as a plain text to the file `tree.txt`

Golem stores its data in `$HOME/.golem/`. The whole tree is stored in `tree.json`. An `aliases` file can be edited to alias command patterns.

## Legacy version

Golem was first developed as an offline browser application. After some regular use I found this too cumbersome so I rewrote it to be a command line tool. You can still find the old html version in `legacy/`.


