package steganography_test

import (
	steganography "camo/steganography"
	"testing"
)


func TestEncryptDecrypt(t *testing.T) {
	data := make([]byte, 100)
	for i:=0; i<100; i++ {
		data[i] = uint8(i)
	}
	encryptor := steganography.Encryptor{}
	encrypted, err := encryptor.Encrypt(data)
	if err != nil {
		t.Errorf("Encoding error: %s", err)
	}
	data_padded := make([]byte, 200)
	_ = copy(data_padded, encrypted)
	for i:=128; i<200; i++ {
		data_padded[i] = 0
	}
	decrypted, err := encryptor.Decrypt(data_padded)
	if err != nil {
		t.Errorf("Decoding error: %s", err)
	}

	for i:=0; i<len(data); i++ {
		if data[i]!=decrypted[i] {
			t.Errorf("Unmatched data at %d", i)
		}
	}
}