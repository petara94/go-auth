package pkg

import "crypto/sha256"

func HashPassword(password string) string {
	return string(sha256.New().Sum([]byte(password)))
}

func PasswordEqual(password, hash string) bool {
	return string(sha256.New().Sum([]byte(password))) == hash
}
