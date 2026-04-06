package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(p string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(p), 10)
	return string(hash)
}
func CheckHash(p, h string) bool {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(p)) == nil
}
