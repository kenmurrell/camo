package steganography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

const key string = "this is like, a key dude"

type Encryptor struct{
}

func (h *Encryptor) Encrypt(in []byte) ([]byte, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }

    return gcm.Seal(nonce, nonce, in, nil), nil
}

func (h *Encryptor) Decrypt(in []byte)  ([]byte, error) { 
    c, err := aes.NewCipher([]byte(key))
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(c)
    if err != nil {
        return nil, err
    }

    nonceSize := gcm.NonceSize()
	in_rm := paddingRemoval(in)
    if len(in_rm) < nonceSize {
        return nil, fmt.Errorf("encrypted data is too short; %d vs %d nonce", len(in_rm), nonceSize)
    }

    nonce, ciph := in_rm[:nonceSize], in_rm[nonceSize:]
    return gcm.Open(nil, nonce, ciph, nil)
}

func paddingRemoval(data []byte) []byte {
	var validIndex int
	for i:=len(data)-1; i>=0; i-- {
		if data[i] != 0 {
			validIndex = i
			break
		}
	}
	return data[:validIndex+1]
}
