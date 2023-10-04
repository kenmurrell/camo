package main

import (
	"fmt"
	"io"
	"os"
)

func printAndExit(err error) {
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

func checkPNG(file *os.File) bool {
	pngsig := [4]byte{0x89, 0x50, 0x4E, 0x47}
	file.Seek(0, io.SeekStart)
	fsig := make([]byte, len(pngsig))
	if _, err := file.Read(fsig); err != nil {
		return false
	}
	file.Seek(0, io.SeekStart)
	//TODO: make reflect.deepequal work here...
	for i:=0; i<len(fsig); i++ {
		if pngsig[i] != fsig[i] {
			return false
		}
	}
	return true
}


