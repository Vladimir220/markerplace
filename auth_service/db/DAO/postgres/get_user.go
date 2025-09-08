package postgres

import (
	"auth_service/models"
	"database/sql"
	"errors"
	"fmt"
)

func (md AuthMarketplaceDAO) GetUser(login string) (user models.User, password string, isFound bool, err error) {
	if login == "" {
		err = errors.New("AuthMarketplaceDAO:GetUser: login not specified")
		return
	}

	queryStr := "SELECT login, group_name, password FROM users WHERE login=$1;"

	err = md.connection.QueryRow(queryStr, login).Scan(&user.Login, &user.Group, &password)
	if err == sql.ErrNoRows {
		err = nil
		return
	} else if err != nil {
		err = fmt.Errorf("AuthMarketplaceDAO:GetUser: %v", err)
		return
	} else {
		isFound = true
	}

	return
}
