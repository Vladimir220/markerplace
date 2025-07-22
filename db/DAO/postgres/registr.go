package postgres

import (
	"errors"
	"fmt"
	"main/models"
)

func (md MarcketplaceDAO) Registr(login, password string) (user models.User, err error) {
	if login == "" || password == "" {
		err = errors.New("MarcketplaceDAO:Registr: login or password not specified")
		return
	}

	queryStr := `INSERT INTO users (login, group_name, password) 
				VALUES ($1, 'user', $2)
				RETURNING login, group_name;`

	connection := md.сonnectionPool.GetConnection()

	err = connection.QueryRow(queryStr, login, password).Scan(&user.Login, &user.Group)
	if err != nil {
		err = fmt.Errorf("MarcketplaceDAO:Registr: %v", err)
		return
	}

	return
}
