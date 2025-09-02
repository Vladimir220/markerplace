package postgres

import (
	"auth_service/models"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type IMarketplaceDAO interface {
	GetUser(login string) (user models.User, password string, isFound bool, err error)
	Registr(login, password string) (user models.User, isAlreadyExist bool, err error)
	Close()
}

func CreateMarketplaceDAO() (dao IMarketplaceDAO, err error) {
	logLabel := "CreateMarketplaceDAO():"

	err = checkDbExistence()
	if err != nil {
		err = fmt.Errorf("%s%v", logLabel, err)
		return
	}

	connection, err := connect()
	if err != nil {
		err = fmt.Errorf("%s%v", logLabel, err)
		return
	}

	err = CheckMigrations(connection)
	if err != nil {
		err = fmt.Errorf("%s%v", logLabel, err)
		return
	}

	mpDao := &MarketplaceDAO{
		connection: connection,
	}

	return mpDao, nil
}

type MarketplaceDAO struct {
	connection *sql.DB
}

func (md MarketplaceDAO) Close() {
	md.connection.Close()
}
