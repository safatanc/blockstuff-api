package converter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func PasswordToHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyPassword(password string, hashPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err
}

func EncryptPassword(password string) (*string, error) {
	block, err := aes.NewCipher([]byte(os.Getenv("API_SECRET_KEY")))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(password), nil)
	encryptedPassword := base64.StdEncoding.EncodeToString(cipherText)

	return &encryptedPassword, nil
}

func DecryptPassword(encryptedPassword string) (*string, error) {
	block, err := aes.NewCipher([]byte(os.Getenv("API_SECRET_KEY")))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	cipherText, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	password, err := gcm.Open(nil, []byte(nonce), []byte(cipherText), nil)
	if err != nil {
		return nil, err
	}

	passwordString := string(password)

	return &passwordString, nil
}
