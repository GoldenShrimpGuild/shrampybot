package utility

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"shrampybot/config"
)

func GenerateRandomHex(byteLength int) string {
	secretKey := make([]byte, byteLength)
	rand.Read(secretKey)
	return hex.EncodeToString(secretKey)
}

func GenerateRandomBase64(byteLength int) string {
	b := make([]byte, byteLength)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// Based on example code from:
// https://bitfieldconsulting.com/posts/aes-encryption

func EncryptSecret(secret string) (string, string, error) {
	key, err := hex.DecodeString(config.DBCryptKey)
	if err != nil {
		return "", "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}

	plaintext := []byte(secret)

	iv := make([]byte, aes.BlockSize)
	_, err = rand.Read(iv)
	if err != nil {
		return "", "", err
	}

	enc := cipher.NewCBCEncrypter(block, iv)
	plaintext = pad(plaintext, aes.BlockSize)
	ciphertext := make([]byte, len(plaintext))
	enc.CryptBlocks(ciphertext, plaintext)

	return hex.EncodeToString(ciphertext), hex.EncodeToString(iv), nil
}

func DecryptSecret(ciphertext string, iv string) (string, error) {
	key, err := hex.DecodeString(config.DBCryptKey)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherbytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	ivbytes, err := hex.DecodeString(iv)
	if err != nil {
		return "", err
	}

	plaintext := make([]byte, len(cipherbytes))
	dec := cipher.NewCBCDecrypter(block, ivbytes)
	dec.CryptBlocks(plaintext, cipherbytes)
	plaintext = unpad(plaintext)

	return string(plaintext), nil
}

func pad(data []byte, blockSize int) []byte {
	n := blockSize - len(data)%blockSize
	padding := bytes.Repeat([]byte{byte(n)}, n)
	return append(data, padding...)
}

func unpad(data []byte) []byte {
	n := int(data[len(data)-1])
	return data[:len(data)-n]
}
