package cipher_client

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func GenerateEncryptionKey() ([]byte, error) {
	key := make([]byte, 32) // 256-bit key for AES-256
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func Encrypt(plainText string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nil, nonce, []byte(plainText), nil)
	encryptedData := append(nonce, cipherText...)
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

func Decrypt(cipherText string, key []byte) (string, error) {
	encryptedData, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, encryptedText := encryptedData[:nonceSize], encryptedData[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, encryptedText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
