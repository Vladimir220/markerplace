package crypto

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GetHashedPassword(password string) (hashedPassword string, err error) {
	logLabel := fmt.Sprintf("GetHashedPassword():[params:%s]:", password)
	res, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}
	hashedPassword = string(res)
	return
}

func ComparePassword(password, hashedPassword string) (equal bool, err error) {
	logLabel := fmt.Sprintf("ComparePassword():[params:%s,%s]:", password, hashedPassword)
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return
	} else if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	} else {
		equal = true
	}
	return
}
