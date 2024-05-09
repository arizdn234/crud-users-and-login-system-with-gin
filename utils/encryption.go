package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(password string) (string, error) {
	hash := sha256.New()

	_, err := hash.Write([]byte(password))
	if err != nil {
		return "", err
	}

	hashedPasswordBytes := hash.Sum(nil)
	hashedPassword := hex.EncodeToString(hashedPasswordBytes)

	return hashedPassword, nil
}
