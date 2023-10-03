package steganography_test

import (
	steganography "camo/steganography"
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"testing"
)

var host1 *os.File
var hide1 *os.File
var cmbn1 *os.File
var rslt1 *os.File

func init() {
	host1, _ = os.Open("test_images/host1.png")
	hide1, _ = os.Open("test_images/hide1.jpg")
	cmbn1, _ = os.Create("test_images/cmbn1.png")
	rslt1, _ = os.Create("test_images/rslt1.jpg")
}
func Shutdown() {
	defer host1.Close()
	defer hide1.Close()
	defer cmbn1.Close()
	defer rslt1.Close()
}

func TestAllRgbaEncoding(t *testing.T) {
	encodeAndDecode(t, steganography.AllRGBA)
}

func TestOnlyBlueEncoding(t *testing.T) {
	encodeAndDecode(t, steganography.BlueRGBA)
}

func encodeAndDecode(t *testing.T, m steganography.Mode) {
	defer Shutdown()
	err := steganography.Encode(host1, hide1, cmbn1, m)
    if err != nil {
        t.Errorf("Error encountered encoding: %s", err.Error())
    }
	cmbn1.Seek(0, io.SeekStart)
	err = steganography.Decode(cmbn1, rslt1, m)
    if err != nil {
        t.Errorf("Error encountered decoding: %s", err.Error())
    }

	//verify
	fi, _ := hide1.Stat()
	hide1.Seek(0, io.SeekStart)
	h1 := md5.New()
	if _, err := io.CopyN(h1, hide1, fi.Size()-1); err != nil {
		log.Fatal(err)
	}
	hash1 := hex.EncodeToString(h1.Sum(nil))
	rslt1.Seek(0, io.SeekStart)
	h2 := md5.New()
	if _, err := io.CopyN(h2, rslt1, fi.Size()-1); err != nil {
		log.Fatal(err)
	}
	hash2 := hex.EncodeToString(h2.Sum(nil))

	if hash1 != hash2 {
		t.Errorf("Images are different! %s vs %s", hash1, hash2)
	}
}