package postgres

import (
	"auth_service/models"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type IAuthMarketplaceDAO interface {
	GetUser(login string) (user models.User, password string, isFound bool, err error)
	Registr(login, password string) (user models.User, isAlreadyExist bool, err error)
	Close()
}

func CreateAuthMarketplaceDAO() (dao IAuthMarketplaceDAO, err error) {
	logLabel := "CreateAuthMarketplaceDAO():"

	err = checkDbExistence(5, time.Second)
	if err != nil {
		err = fmt.Errorf("%s%v", logLabel, err)
		return
	}

	connection, err := connect()
	if err != nil {
		err = fmt.Errorf("%s%v", logLabel, err)
		return
	}

	mpDao := &AuthMarketplaceDAO{
		connection: connection,
	}

	return mpDao, nil
}

type AuthMarketplaceDAO struct {
	connection *sql.DB
}

func (md AuthMarketplaceDAO) Close() {
	md.connection.Close()
}
