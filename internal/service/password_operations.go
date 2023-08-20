package service

import (
	"crypto/sha1"
	"fmt"
)

func generatePasswordHash(password, secret string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(secret)))
}

func comparePasswordHash(hashedPassword, password, secret string) bool {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(secret))) == hashedPassword
}
