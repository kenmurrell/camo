package main

import (
	"flag"
	"fmt"
	"os"
)

func run() {
	encodeCmd := flag.NewFlagSet("encode", flag.ExitOnError)
	hostPtr := encodeCmd.String("host", "", "The image to host a hidden file.")
	hidePtr := encodeCmd.String("hide", "", "The image to hide.")

	decodeCmd := flag.NewFlagSet("decode", flag.ExitOnError)
	hostedPtr := decodeCmd.String("host", "", "The image that contains a hidden file.")
	outputPtr := decodeCmd.String("output", "", "The name of the output file.")

	switch os.Args[1] {
	case "encode":
		encodeCmd.Parse(os.Args[2:])
		encodeWork(*hostPtr, *hidePtr)
	case "decode":
		decodeCmd.Parse(os.Args[2:])
		decodeWork(*hostedPtr, *outputPtr)
	default:
		fmt.Println("Use \"encode\" to encode images, \"decode\" to decode them.\nRun \"encode -h\" to show flags.")
		os.Exit(1)
	}
}

func encodeWork(hostFileName string, hideFileName string) {
	hostFile, err := loadFile(hostFileName)
	defer hostFile.Close()
	_printAndExit(err)
	hideFile, err := loadFile(hideFileName)
	defer hideFile.Close()
	_printAndExit(err)
	hostedFileName := hostFileName[:len(hostFileName)-4] + "_hosted.png"
	hostedFile, err := os.Create(hostedFileName)
	defer hostedFile.Close()
	_printAndExit(err)

	err = encode(hostFile, hideFile, hostedFile)
	_printAndExit(err)
}

func decodeWork(hostFileName string, hideFileName string) {
	hostFile, err := loadFile(hostFileName)
	_printAndExit(err)
	defer hostFile.Close()
	err = decode(hostFile, hideFileName)
	_printAndExit(err)
}
