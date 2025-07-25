package postgres

import (
	"errors"
	"fmt"
	"main/models"
)

func (md MarcketplaceDAO) GetUser(login string) (user models.User, password string, isFound bool, err error) {
	if login == "" {
		err = errors.New("MarcketplaceDAO:GetUser: login not specified")
		return
	}

	queryStr := "SELECT login, group_name, password FROM users WHERE login=$1;"

	connection := md.сonnectionPool.GetConnection()

	err = connection.QueryRow(queryStr, login).Scan(&user.Login, &user.Group, &password)
	if err != nil {
		err = fmt.Errorf("MarcketplaceDAO:GetUser: %v", err)
		return
	}

	isFound = true
	return
}
