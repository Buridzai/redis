package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// HashPasswordNoErr tiện cho seed dữ liệu (tự động panic nếu lỗi)
func HashPasswordNoErr(password string) string {
	hashed, err := HashPassword(password)
	if err != nil {
		panic(err)
	}
	return hashed
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
