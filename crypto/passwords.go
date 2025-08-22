package crypto

import "golang.org/x/crypto/bcrypt"

func GetHashedPassword(password string) (hashedPassword string, err error) {
	res, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	hashedPassword = string(res)
	return
}

func ComparePassword(password, hashedPassword string) (equal bool) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == nil {
		equal = true
	}
	return
}
