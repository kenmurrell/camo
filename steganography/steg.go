package steganography

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
)

type nrgbaEncoder struct {
	input *image.NRGBA
}

func (e *nrgbaEncoder) encode(data []byte) (*image.NRGBA, int) {
	bounds := (*e.input).Bounds()
	newImg := image.NewNRGBA(bounds)
	bitCount := 0
	byteCount := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			nrgba := (*e.input).NRGBAAt(x, y)
			b := ((nrgba.B >> 1) << 1)
			if byteCount < len(data) {
				b += uint8((data[byteCount] >> (bitCount % 8)) & 0x1)
				bitCount++
				if bitCount % 8 == 0 {
					byteCount++
				}
			}
			newImg.SetNRGBA(x, y, color.NRGBA{
				R: nrgba.R,
				G: nrgba.G,
				B: b,
				A: nrgba.A,
			})
		}
	}
	return newImg, byteCount
}

func (e *nrgbaEncoder) decode(data []byte) int {
	bounds := (*e.input).Bounds()
	bitCount := 0
	byteCount := 0
	var cnstrByte byte = 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			nrgba := (*e.input).NRGBAAt(x, y)
			cnstrByte += byte((nrgba.B & 0x1) << (bitCount % 8))
			bitCount++
			if bitCount % 8 == 0 {
				data[byteCount] = cnstrByte
				cnstrByte = 0
				byteCount++
			}
		}
	}
	return byteCount
}

func Encode(hostFile *os.File, hideFile *os.File, hostedFile *os.File) error {
	fmt.Print("Encoding...")
	hostIm, _, err := image.Decode(hostFile)
	if err != nil {
		return fmt.Errorf("ERROR: image cannot be decoded %s", err.Error())
	}
	fi, _ := hideFile.Stat()
	data := make([]byte, fi.Size())
	hideFile.Read(data)
	bounds := hostIm.Bounds()

	capacity := int((bounds.Max.Y - bounds.Min.Y) * (bounds.Max.X - bounds.Min.X) / 8)
	if len(data) > capacity {
		return fmt.Errorf("ERROR: ghost file exceeds host file's capacity (>%d bytes)", capacity)
	}
	if hostImNRGBA, ok := hostIm.(*image.NRGBA); ok {
		nrgbaEncoder := nrgbaEncoder{hostImNRGBA}
		newImg, nencoded := nrgbaEncoder.encode(data)
		if(nencoded < len(data)) {
			fmt.Printf("WARN: only %d bytes encoded", nencoded)
		}
		if err := png.Encode(hostedFile, newImg); err != nil {
			return err
		}
	}
	fmt.Println("...Done.")
	return nil
}

func Decode(hostFile *os.File, outputFile *os.File) error {
	fmt.Print("Decoding...")
	hostIm, _, err := image.Decode(hostFile)
	if err != nil {
		return fmt.Errorf("ERROR: image cannot be decoded %s", err.Error())
	}
	bounds := hostIm.Bounds()
	capacity := int((bounds.Max.Y - bounds.Min.Y) * (bounds.Max.X - bounds.Min.X) / 8)
	data := make([]byte, capacity)

	if hostImNRGBA, ok := hostIm.(*image.NRGBA); ok {
		nrgbaEncoder := nrgbaEncoder{hostImNRGBA}
		_ = nrgbaEncoder.decode(data)
	}
	if _, err = outputFile.Write(data); err != nil {
		return err
	}

	fmt.Println("...Done.")
	return nil
}