package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"main/models"
)

func (md MarketplaceDAO) Registr(login, password string) (user models.User, isAlreadyExist bool, err error) {
	if login == "" || password == "" {
		err = errors.New("MarketplaceDAO:Registr: login or password not specified")
		return
	}

	queryStr := `INSERT INTO users (login, group_name, password) 
				VALUES ($1, 'user', $2)
				ON CONFLICT (login) DO NOTHING
				RETURNING login, group_name;`

	connection := md.сonnectionPool.GetConnection()

	err = connection.QueryRow(queryStr, login, password).Scan(&user.Login, &user.Group)
	if err == sql.ErrNoRows {
		isAlreadyExist = true
	} else if err != nil {
		err = fmt.Errorf("MarketplaceDAO:Registr: %v", err)
		return
	}

	return
}
