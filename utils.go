package main

import (
	"fmt"
	"os"
)

func _printAndExit(err error) {
	if err != nil {
		fmt.Println("\n" + err.Error())
		os.Exit(1)
	}
}

func loadFile(filepath string) (*os.File, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("ERROR: file %s not found", filepath)
	}
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("ERROR: file %s cannot be read", filepath)
	}
	return file, nil
}
