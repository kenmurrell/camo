package steganography_test

import (
	steganography "camo/steganography"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
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

func TestEncryptionEncoding(t *testing.T) {
	defer Shutdown()
	r := steganography.RunOptions{
		Mode: steganography.AllRGBA,
		Encrypt: true,
	}
	_ = encodeAndDecode(r)
	compareFileData(t, hide1, rslt1)
	compareFileHashes(t, hide1, rslt1)
}

func TestAllRgbaEncoding(t *testing.T) {
	defer Shutdown()
	r := steganography.RunOptions{
		Mode: steganography.AllRGBA,
		Encrypt: false,
	}
	_ = encodeAndDecode(r)
	compareFileHashes(t, hide1, rslt1)
}

func TestOnlyBlueEncoding(t *testing.T) {
	defer Shutdown()
	r := steganography.RunOptions{
		Mode: steganography.BlueRGBA,
		Encrypt: false,
	}
	_ = encodeAndDecode(r)
	compareFileHashes(t, hide1, rslt1)
}

func encodeAndDecode(r steganography.RunOptions) error {
	err := steganography.Encode(host1, hide1, cmbn1, r)
    if err != nil {
        return fmt.Errorf("Error encountered encoding: %s", err.Error())
    }
	cmbn1.Seek(0, io.SeekStart)
	err = steganography.Decode(cmbn1, rslt1, r)
    if err != nil {
        return fmt.Errorf("Error encountered decoding: %s", err.Error())
    }
	return nil
}

func compareFileHashes(t *testing.T, file1 *os.File, file2 *os.File){
	fi, _ := file1.Stat()
	file1.Seek(0, io.SeekStart)
	h1 := md5.New()
	if _, err := io.CopyN(h1, file1, fi.Size()-1); err != nil {
		t.Error("Unable to hash file1")
	}
	hash1 := hex.EncodeToString(h1.Sum(nil))
	file2.Seek(0, io.SeekStart)
	h2 := md5.New()
	if _, err := io.CopyN(h2, file2, fi.Size()-1); err != nil {
		t.Error("Unable to hash file2")
	}
	hash2 := hex.EncodeToString(h2.Sum(nil))

	if hash1 != hash2 {
		t.Errorf("Images are different! %s vs %s", hash1, hash2)
	}
}

func compareFileData(t *testing.T, file1 *os.File, file2 *os.File){
	fi, _ := file1.Stat()
	file1.Seek(0, io.SeekStart)
	file1Arr := make([]byte, fi.Size())
	_, _ = file1.Read(file1Arr)
	file2.Seek(0, io.SeekStart)
	file2Arr := make([]byte, fi.Size())
	_, _ = file2.Read(file2Arr)
	for i:=0; i<len(file2Arr); i++ {
		if file1Arr[i]!=file2Arr[i] {
			t.Errorf("Unmatched data at %d", i)
			break;
		}
	}
}