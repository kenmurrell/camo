package main

import (
	stg "camo/steganography"
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
		encodeWork(*hostPtr, *hidePtr, buildRunOptions(eModePtr, eEncrypt))
	case "decode":
		decodeCmd.Parse(os.Args[2:])
		decodeWork(*hostedPtr, *outputPtr, buildRunOptions(dModePtr, dEncrypt))
	default:
		fmt.Println("Use \"encode\" to encode images, \"decode\" to decode them.\nRun \"encode -h\" to show flags.")
		os.Exit(1)
	}
}

func encodeWork(hostFileName string, hideFileName string, ro stg.RunOptions) {
	hostFile, err := loadFile(hostFileName)
	printAndExit(err)
	defer hostFile.Close()
	if !checkPNG(hostFile) {
		fmt.Println("Only PNG files are supported as the host file!")
		os.Exit(1)
	}
	hideFile, err := loadFile(hideFileName)
	printAndExit(err)
	defer hideFile.Close()
	hostedFileName := hostFileName[:len(hostFileName)-4] + "_hosted.png"
	hostedFile, err := os.Create(hostedFileName)
	printAndExit(err)
	defer hostedFile.Close()
	ro.Print()
	err = stg.Encode(hostFile, hideFile, hostedFile, ro)
	printAndExit(err)
}

func decodeWork(hostFileName string, rsltFileName string, ro stg.RunOptions) {
	hostFile, err := loadFile(hostFileName)
	printAndExit(err)
	defer hostFile.Close()
	rsltFile, err := os.Create(rsltFileName)
	printAndExit(err)
	defer rsltFile.Close()
	ro.Print()
	err = stg.Decode(hostFile, rsltFile, ro)
	printAndExit(err)
}

func buildRunOptions(eModePtr *bool, encrypt *bool) stg.RunOptions {
	mode := stg.AllRGBA
	if *eModePtr { 
		mode = stg.BlueRGBA
	}
	return stg.RunOptions{ Mode: mode, Encrypt: *encrypt }
}
