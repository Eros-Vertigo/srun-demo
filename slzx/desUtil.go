package slzx

import (
	"bytes"
	"crypto/des"
	"encoding/base64"
)

func pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

func unPad(src []byte) []byte {
	length := len(src)
	unPadding := int(src[length-1])
	return src[:(length - unPadding)]
}

// EncryptDES 加密
func EncryptDES(key, plaintext string) (string, error) {
	// Convert strings to bytes
	keyBytes := []byte(key)
	plaintextBytes := []byte(plaintext)

	// Create DES block cipher
	block, err := des.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Pad plaintext to be a multiple of the block size
	paddedText := pad(plaintextBytes, block.BlockSize())

	// Create ciphertext array
	ciphertext := make([]byte, len(paddedText))

	// Encrypt each block
	for bs, be := 0, block.BlockSize(); bs < len(paddedText); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Encrypt(ciphertext[bs:be], paddedText[bs:be])
	}

	// Encode ciphertext to base64
	encodedText := base64.StdEncoding.EncodeToString(ciphertext)
	return encodedText, nil
}

// DecryptDES 解密
func DecryptDES(key, encryptedText string) (string, error) {
	keyBytes := []byte(key)
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := des.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	decrypted := make([]byte, len(ciphertext))

	for bs, be := 0, block.BlockSize(); bs < len(ciphertext); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Decrypt(decrypted[bs:be], ciphertext[bs:be])
	}

	unPaddedText := unPad(decrypted)
	return string(unPaddedText), nil
}
