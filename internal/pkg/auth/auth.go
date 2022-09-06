package auth

import "golang.org/x/crypto/bcrypt"

// Encrypt 对字符串加密（hash）
func Encrypt(source string) (string, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(source), bcrypt.DefaultCost)
	return string(hashedPwd), err
}

// Compare 比较加密前、后的密码是否一致
func Compare(hashedPwd, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(password))
}
