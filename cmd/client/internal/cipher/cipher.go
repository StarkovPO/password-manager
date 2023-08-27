package cipher_client

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

func GetUserPath() (string, error) {
	keyFilePath := ""
	homeDir, err := os.UserHomeDir()

	if err != nil {
		fmt.Println("Error getting user's home directory:", err)
		return "", err
	}

	switch runtime.GOOS {
	case "windows":
		keyFilePath = filepath.Join(homeDir, "AppData", "Roaming", "password-saver", "encryption.json")
	case "darwin":
		keyFilePath = filepath.Join(homeDir, "password-saver", "encryption.json")
	default:
		keyFilePath = filepath.Join(homeDir, "password-saver", "encryption.json")
	}

	if exist := checkDirAndCreate(keyFilePath); !exist {
		fmt.Errorf("unhandler error while create the file")
		return "", errors.New("error while create the file")
	}

	return keyFilePath, nil
}

func checkDirAndCreate(path string) bool {

	p := strings.TrimSuffix(path, "/encryption.json")

	_, err := os.Stat(p)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		err = os.Mkdir(p, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return false
		}
	}
	return true
}

func SaveEncryptedKey(userID string, key []byte) error {

	f, err := GetUserPath()
	if err != nil {
		return err
	}

	p, err := NewProducer(f)
	if err != nil {
		return err
	}

	err = p.WriteFile(userID, key)
	if err != nil {
		return err
	}

	return nil
}

func ReadEncryptedKey() (map[string]string, error) {
	UserKeys := make(map[string]string)
	f, err := GetUserPath()

	if err != nil {
		return nil, err
	}
	c, err := NewConsumer(f)

	if err != nil {
		return nil, err
	}

	type tmp struct {
		UID string `json:"UID"`
		Key string `json:"Key"`
	}

	b, err := c.ReadFile()

	if err != nil {
		return nil, err
	}

	var result []tmp

	err = json.Unmarshal(b, &result)

	if err != nil {
		return nil, err
	}
	for _, v := range result {
		UserKeys[v.UID] = v.Key
	}

	return UserKeys, nil
}
