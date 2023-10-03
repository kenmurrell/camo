package main

import (
	steganography "camo/steganography"
	"flag"
	"fmt"
	"os"
)

func run() {
	encodeCmd := flag.NewFlagSet("encode", flag.ExitOnError)
	eModePtr := encodeCmd.Bool("blue", false, "Encode only in blue pixels.")
	eEncrypt := encodeCmd.Bool("encrypt", false, "Encrypt the encoded content.")
	hostPtr := encodeCmd.String("host", "", "The image to host a hidden file.")
	hidePtr := encodeCmd.String("hide", "", "The image to hide.")

	decodeCmd := flag.NewFlagSet("decode", flag.ExitOnError)
	dModePtr := decodeCmd.Bool("blue", false, "Decode only the blue pixels.")
	dEncrypt := decodeCmd.Bool("decrypt", false, "Decrypt the encoded content.")
	hostedPtr := decodeCmd.String("host", "", "The image that contains a hidden file.")
	outputPtr := decodeCmd.String("output", "", "The name of the output file.")

	switch os.Args[1] {
	case "encode":
		encodeCmd.Parse(os.Args[2:])
		encodeWork(*hostPtr, *hidePtr, *eModePtr, *eEncrypt)
	case "decode":
		decodeCmd.Parse(os.Args[2:])
		decodeWork(*hostedPtr, *outputPtr, *dModePtr, *dEncrypt)
	default:
		fmt.Println("Use \"encode\" to encode images, \"decode\" to decode them.\nRun \"encode -h\" to show flags.")
		os.Exit(1)
	}
}

func encodeWork(hostFileName string, hideFileName string, modeB bool, encr bool) {
	hostFile, err := loadFile(hostFileName)
	_printAndExit(err)
	defer hostFile.Close()
	hideFile, err := loadFile(hideFileName)
	_printAndExit(err)
	defer hideFile.Close()
	hostedFileName := hostFileName[:len(hostFileName)-4] + "_hosted.png"
	hostedFile, err := os.Create(hostedFileName)
	_printAndExit(err)
	defer hostedFile.Close()
	mode := steganography.AllRGBA
	if modeB { 
		mode = steganography.BlueRGBA
	}
	r := steganography.RunOptions{
		Mode: mode,
		Encrypt: encr,
	}
	err = steganography.Encode(hostFile, hideFile, hostedFile, r)
	_printAndExit(err)
}

func decodeWork(hostFileName string, rsltFileName string, modeB bool, encr bool) {
	hostFile, err := loadFile(hostFileName)
	_printAndExit(err)
	defer hostFile.Close()
	rsltFile, err := os.Create(rsltFileName)
	_printAndExit(err)
	defer rsltFile.Close()
	mode := steganography.AllRGBA
	if modeB { 
		mode = steganography.BlueRGBA
	}
	r := steganography.RunOptions{
		Mode: mode,
		Encrypt: encr,
	}
	err = steganography.Decode(hostFile, rsltFile, r)
	_printAndExit(err)
}
