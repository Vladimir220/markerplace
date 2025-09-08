package postgres

import (
	"auth_service/models"
	"database/sql"
	"errors"
	"fmt"
)

func (md AuthMarketplaceDAO) Registr(login, password string) (user models.User, isAlreadyExist bool, err error) {
	if login == "" || password == "" {
		err = errors.New("AuthMarketplaceDAO:Registr: login or password not specified")
		return
	}

	queryStr := `INSERT INTO users (login, group_name, password) 
				VALUES ($1, 'user', $2)
				ON CONFLICT (login) DO NOTHING
				RETURNING login, group_name;`

	err = md.connection.QueryRow(queryStr, login, password).Scan(&user.Login, &user.Group)
	if err == sql.ErrNoRows {
		isAlreadyExist = true
		err = nil
	} else if err != nil {
		err = fmt.Errorf("AuthMarketplaceDAO:Registr: %v", err)
		return
	}

	return
}
