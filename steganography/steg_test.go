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

const host1filename string = "test_images/host1.png"
const hide1filename string = "test_images/hide1.jpg"
const cmbn1filename string = "test_images/cmbn1.png"
const rslt1filename string = "test_images/rslt1.jpg"

func TestEncodeAndDecode(t *testing.T) {
	host1, _ := os.Open(host1filename)
	defer host1.Close()
	hide1, _ := os.Open(hide1filename)
	defer hide1.Close()
	cmbn1, _ := os.Create(cmbn1filename)
	defer cmbn1.Close()
	rslt1, _ := os.Create(rslt1filename)
	defer rslt1.Close()

	err := steganography.Encode(host1, hide1, cmbn1)
    if err != nil {
        t.Errorf("Error encountered encoding: %s", err.Error())
    }
	cmbn1.Seek(0, io.SeekStart)
	err = steganography.Decode(cmbn1, rslt1)
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