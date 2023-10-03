package steganography

import (
	"image"
	"image/color"
)

type nrbgaBlueEncoder struct {
	input *image.NRGBA
}

func (e *nrbgaBlueEncoder) encode(data []byte) (*image.NRGBA, int) {
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
				if bitCount%8 == 0 {
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

func (e *nrbgaBlueEncoder) decode(data []byte) int {
	bounds := (*e.input).Bounds()
	bitCount := 0
	byteCount := 0
	var cnstrByte byte = 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			nrgba := (*e.input).NRGBAAt(x, y)
			cnstrByte += byte((nrgba.B & 0x1) << (bitCount % 8))
			bitCount++
			if bitCount%8 == 0 {
				data[byteCount] = cnstrByte
				cnstrByte = 0
				byteCount++
			}
		}
	}
	return byteCount
}

type nrbgaAllEncoder struct {
	input *image.NRGBA
}

func (e *nrbgaAllEncoder) encode(data []byte) (*image.NRGBA, int) {
	bounds := (*e.input).Bounds()
	newImg := image.NewNRGBA(bounds)
	bitCount := 0
	byteCount := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			nrgba := (*e.input).NRGBAAt(x, y)
			r := ((nrgba.R >> 1) << 1)
			g := ((nrgba.G >> 1) << 1)
			b := ((nrgba.B >> 1) << 1)
			a := ((nrgba.A >> 1) << 1)
			if byteCount < len(data) {
				r += uint8((data[byteCount] >> (bitCount % 8)) & 0x1)
				g += uint8((data[byteCount] >> ((bitCount + 1) % 8)) & 0x1)
				b += uint8((data[byteCount] >> ((bitCount + 2) % 8)) & 0x1)
				a += uint8((data[byteCount] >> ((bitCount + 3) % 8)) & 0x1)
				bitCount+=4
				if bitCount%8 == 0 {
					byteCount++
				}
			}
			newImg.SetNRGBA(x, y, color.NRGBA{
				R: r,
				G: g,
				B: b,
				A: a,
			})
		}
	}
	return newImg, byteCount
}

func (e *nrbgaAllEncoder) decode(data []byte) int {
	bounds := (*e.input).Bounds()
	bitCount := 0
	byteCount := 0
	var cnstrByte byte = 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			nrgba := (*e.input).NRGBAAt(x, y)
			cnstrByte += byte((nrgba.R & 0x1) << (bitCount % 8))
			cnstrByte += byte((nrgba.G & 0x1) << ((bitCount + 1) % 8))
			cnstrByte += byte((nrgba.B & 0x1) << ((bitCount + 2) % 8))
			cnstrByte += byte((nrgba.A & 0x1) << ((bitCount + 3) % 8))
			bitCount+=4
			if bitCount%8 == 0 {
				data[byteCount] = cnstrByte
				cnstrByte = 0
				byteCount++
			}
		}
	}
	return byteCount
}