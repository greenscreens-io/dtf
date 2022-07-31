package model

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

// EncryptMessage encrypt string into base64
func EncryptMessage(key []byte, message string) (string, error) {
	cipherText, err := EncryptMessageRaw(key, []byte(message))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// EncryptMessageRaw encrypt bytes to bytes
func EncryptMessageRaw(key []byte, message []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create new cipher: %v", err)
	}

	cipherText := make([]byte, aes.BlockSize+len(message))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("could not encrypt: %v", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], message)

	return cipherText, nil
}

// DecryptMessage decrpt bas64 to string
func DecryptMessage(key []byte, message string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", fmt.Errorf("could not base64 decode: %v", err)
	}

	cipherText, err = DecryptMessageRaw(key, cipherText)
	if err != nil {
		return "", err
	}

	return string(cipherText), nil
}

// DecryptMessageRaw decrypt bytes to bytes
func DecryptMessageRaw(key []byte, message []byte) ([]byte, error) {
	cipherText := message

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create new cipher: %v", err)
	}

	if len(cipherText) < aes.BlockSize {
		return nil, fmt.Errorf("invalid ciphertext block size")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, nil
}

func Protect(appID, id string) []byte {
	mac := hmac.New(sha256.New, []byte(id))
	mac.Write([]byte(appID))
	return mac.Sum(nil)[0:16]
}
