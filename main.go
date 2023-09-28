package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
)

const loadedhost string = "loaded_host.png"
const decodedhidden string = "decode.png"

func loadFile(filepath string) ([]byte, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("ERROR: file %s not found", filepath)
	}
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("ERROR: file %s cannot be read", filepath)
	}
	return content, nil
}

func loadImage(filepath string) (*image.Image, error) {
	//TODO: check if file is png/jpg/gif
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("ERROR: file %s not found", filepath)
	}
	reader, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("ERROR: file %s cannot be read", filepath)
	}
	defer reader.Close()

	im, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	return &im, nil
}

func encode(ghostBytes *[]byte, hostIm *image.Image) error {
	fmt.Println("Encoding..")
	bounds := (*hostIm).Bounds()
	capacity := int((bounds.Max.Y - bounds.Min.Y) * (bounds.Max.X - bounds.Min.X) / 8)
	if len(*ghostBytes) > capacity {
		return fmt.Errorf("ERROR: ghost file exceeds host file's capacity (>%d bytes)", capacity)
	}

	totalCtr := 0
	bitPos := 0
	c := 0
	newImg := image.NewGray16(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			nrgba, _ := (*hostIm).At(x, y).(color.NRGBA)
			alpha, _ := color.Gray16Model.Convert(nrgba).(color.Gray16)
			b := alpha.Y

			// remove last bit
			b = ((b >> 1) << 1)
			if c < len(*ghostBytes) {
				k := (*ghostBytes)[c]
				bit := (k >> bitPos) & 0x1
				b += uint16(bit)
			}

			newImg.SetGray16(x, y, color.Gray16{b})
			totalCtr++
			bitPos = totalCtr % 8
			if bitPos == 0 {
				c++
			}
		}
	}

	f, err := os.Create(loadedhost)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := png.Encode(f, newImg); err != nil {
		f.Close()
		return err
	}

	return nil
}

func decode(hostIm *image.Image) error {
	fmt.Println("Decoding..")
	bounds := (*hostIm).Bounds()
	capacity := int((bounds.Max.Y - bounds.Min.Y) * (bounds.Max.X - bounds.Min.X) / 8)
	ghost := make([]byte, capacity)

	totalCtr := 0
	bytePos := 0
	c := 0
	var myByte byte = 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			alpha, _ := color.Gray16Model.Convert((*hostIm).At(x, y)).(color.Gray16)
			b := alpha.Y
			myByte += byte((b & 0x1) << bytePos)

			totalCtr++
			bytePos = totalCtr % 8
			if bytePos == 0 {
				ghost[c] = myByte
				myByte = 0
				c++
			}
		}
	}

	err := os.WriteFile(decodedhidden, ghost, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func main() {
	ghost := "images/tree-small.png"
	host := "images/honda.png"
	if host[len(host)-3:] != "png" {
		fmt.Printf("ERROR: only supports png host file.")
		os.Exit(1)
	}

	ghostBytes, err := loadFile(ghost)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	hostImPtr, err := loadImage(host)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	hostIm := *hostImPtr
	//TODO: check the stenographical capacity of the file
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = encode(&ghostBytes, &hostIm)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	hostImPtr2, err := loadImage(loadedhost)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	hostIm2 := *hostImPtr2
	//TODO: check the stenographical capacity of the file
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = decode(&hostIm2)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
