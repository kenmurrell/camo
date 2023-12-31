package steganography

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
)

type RunOptions struct {
	Mode    EncodingMode
	Encrypt bool
}

type EncodingMode byte

const (
	BlueRGBA EncodingMode = iota
	AllRGBA  EncodingMode = iota
)

func (ro *RunOptions) Print() {
	encryptStr := "OFF"
	if ro.Encrypt {
		encryptStr = "ON"
	}
	var modeStr string
	switch ro.Mode {
	case 0x0:
		modeStr = "BlueRGBA"
	case 0x1:
		modeStr = "AllRGBA"
	default:
		modeStr = "none?"
	}
	fmt.Printf("Mode:%s Encryption:%s\n", modeStr, encryptStr)
}

func Encode(hostFile *os.File, hideFile *os.File, hostedFile *os.File, r RunOptions) error {
	fmt.Print("Encoding...")
	hostIm, _, err := image.Decode(hostFile)
	if err != nil {
		return fmt.Errorf("ERROR: image cannot be decoded %s", err.Error())
	}
	fi, _ := hideFile.Stat()
	data := make([]byte, fi.Size())
	hideFile.Read(data)
	if r.Encrypt {
		encryptor := Encryptor{}
		if data, err = encryptor.Encrypt(data); err != nil {
			return fmt.Errorf("ERROR: image cannot be encrypted %s", err.Error())
		}
	}
	bounds := hostIm.Bounds()
	capacity := findCapacity(&bounds, r.Mode)
	if len(data) > capacity {
		return fmt.Errorf("ERROR: file exceeds host file's encoding capacity (>%d bytes)", capacity)
	}
	if hostImNRGBA, ok := hostIm.(*image.NRGBA); ok {
		var newImg *image.NRGBA
		var nencoded int
		if r.Mode == AllRGBA {
			encoder := nrbgaAllEncoder{hostImNRGBA}
			newImg, nencoded = encoder.encode(data)
		} else {
			encoder := nrbgaBlueEncoder{hostImNRGBA}
			newImg, nencoded = encoder.encode(data)
		}

		if nencoded < len(data) {
			fmt.Printf("WARN: only %d bytes encoded", nencoded)
		}
		if err := png.Encode(hostedFile, newImg); err != nil {
			return err
		}
	}
	fmt.Println("...Done.")
	return nil
}

func Decode(hostFile *os.File, outputFile *os.File, r RunOptions) error {
	fmt.Print("Decoding...")
	hostIm, _, err := image.Decode(hostFile)
	if err != nil {
		return fmt.Errorf("ERROR: image cannot be decoded %s", err.Error())
	}
	bounds := hostIm.Bounds()
	capacity := int((bounds.Max.Y - bounds.Min.Y) * (bounds.Max.X - bounds.Min.X) / 2)
	data := make([]byte, capacity)
	if hostImNRGBA, ok := hostIm.(*image.NRGBA); ok {
		if r.Mode == AllRGBA {
			encoder := nrbgaAllEncoder{hostImNRGBA}
			_ = encoder.decode(data)
		} else {
			encoder := nrbgaBlueEncoder{hostImNRGBA}
			_ = encoder.decode(data)
		}
	}
	if r.Encrypt {
		encryptor := Encryptor{}
		if data, err = encryptor.Decrypt(data); err != nil {
			return fmt.Errorf("ERROR: image cannot be decrypted %s", err.Error())
		}
	}
	if _, err = outputFile.Write(data); err != nil {
		return err
	}

	fmt.Println("...Done.")
	return nil
}

func findCapacity(rect *image.Rectangle, mode EncodingMode) int {
	pixelcount := (rect.Max.Y - rect.Min.Y) * (rect.Max.X - rect.Min.X)
	if mode == BlueRGBA {
		return int(pixelcount / 8)
	} else if mode == AllRGBA {
		return int(pixelcount / 2)
	} else {
		return 0
	}
}
