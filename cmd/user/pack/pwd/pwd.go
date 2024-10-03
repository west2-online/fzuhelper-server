package pwd

import (
	"golang.org/x/crypto/bcrypt"
)

// 加密
func SetPassword(password string) (hashBytes string) {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashBytes = string(bytes)
	return hashBytes
}

// 解密
func CheckPassword(pwdDigest string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pwdDigest), []byte(password))
	return err == nil
}
