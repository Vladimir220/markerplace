package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"main/models"
)

func (md MarketplaceDAO) GetUser(login string) (user models.User, password string, isFound bool, err error) {
	if login == "" {
		err = errors.New("MarketplaceDAO:GetUser: login not specified")
		return
	}

	queryStr := "SELECT login, group_name, password FROM users WHERE login=$1;"

	connection := md.сonnectionPool.GetConnection()

	err = connection.QueryRow(queryStr, login).Scan(&user.Login, &user.Group, &password)
	if err == sql.ErrNoRows {
		return
	} else if err != nil {
		err = fmt.Errorf("MarketplaceDAO:GetUser: %v", err)
		return
	} else {
		isFound = true
	}

	return
}
