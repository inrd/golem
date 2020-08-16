package main

import (
	"os"
	"path/filepath"
)

const golemDirName = ".golem"
const golemTreeFileName = "tree.json"

// Returns the Golem directory path
func InitGolemDir() string {
	p := filepath.Join(os.Getenv("HOME"), golemDirName)
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		err = os.Mkdir(p, 0755)
		if err != nil {
			PrintError(err.Error())
		}
	}
	return p
}

func TreeFilePath(golemDirPath string) string {
	return filepath.Join(golemDirPath, golemTreeFileName)
}
