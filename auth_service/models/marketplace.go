package models

type User struct {
	Login string `redis:"login"`
	Group string `redis:"group"`
}
