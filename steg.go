package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
)

func encode(hostFile *os.File, hideFile *os.File, hostedFile *os.File) error {
	fmt.Print("Encoding..")
	hostIm, _, err := image.Decode(hostFile)
	if err != nil {
		return fmt.Errorf("ERROR: image cannot be decoded.")
	}
	fi, _ := hideFile.Stat()
	hideBytes := make([]byte, fi.Size())
	hideFile.Read(hideBytes)
	bounds := hostIm.Bounds()

	capacity := int((bounds.Max.Y - bounds.Min.Y) * (bounds.Max.X - bounds.Min.X) / 8)
	if len(hideBytes) > capacity {
		return fmt.Errorf("ERROR: ghost file exceeds host file's capacity (>%d bytes)", capacity)
	}

	totalCtr := 0
	bitPos := 0
	c := 0
	newImg := image.NewNRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			nrgba, _ := hostIm.At(x, y).(color.NRGBA)
			// remove last bit
			b := nrgba.B
			b = ((b >> 1) << 1)
			if c < len(hideBytes) {
				k := hideBytes[c]
				bit := (k >> bitPos) & 0x1
				b += uint8(bit)
			}

			newImg.SetNRGBA(x, y, color.NRGBA{
				R: nrgba.R,
				G: nrgba.G,
				B: b,
				A: nrgba.A,
			})
			totalCtr++
			bitPos = totalCtr % 8
			if bitPos == 0 {
				c++
			}
		}
	}

	if err := png.Encode(hostedFile, newImg); err != nil {
		return err
	}
	fmt.Println("...Done.")
	return nil
}

func decode(hostFile *os.File, hideFileName string) error {
	fmt.Print("Decoding..")
	hostIm, _, err := image.Decode(hostFile)
	if err != nil {
		return fmt.Errorf("ERROR: image cannot be decoded.")
	}
	bounds := hostIm.Bounds()
	capacity := int((bounds.Max.Y - bounds.Min.Y) * (bounds.Max.X - bounds.Min.X) / 8)
	hideArr := make([]byte, capacity)

	totalCtr := 0
	bytePos := 0
	c := 0
	var myByte byte = 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			nrgba, _ := hostIm.At(x, y).(color.NRGBA)
			b := nrgba.B
			myByte += byte((b & 0x1) << bytePos)

			totalCtr++
			bytePos = totalCtr % 8
			if bytePos == 0 {
				hideArr[c] = myByte
				myByte = 0
				c++
			}
		}
	}

	err = os.WriteFile(hideFileName, hideArr, 0666)
	if err != nil {
		return err
	}

	fmt.Println("...Done.")
	return nil
}
