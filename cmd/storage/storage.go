package storage

import (
	"database/sql"
	"flag"
	"time"
)

var (
	DB         Storage
	postgreSQL = flag.Bool("postgreSQL", false, "")
	mongoDB    = flag.Bool("mongoDB", false, "")
)

func ConnToDB() {
	flag.Parse()
	if *postgreSQL {
		DB = NewPostgreSQL()
		DB.ConnToDB()
	}
	if *mongoDB {
		DB = NewMongoDB()
		DB.ConnToDB()
	}
}

type Storage interface {
	ConnToDB()
	UserExist(string, interface{}) error
	EmailExist(string, interface{}) error
	SaveUser(interface{}) error
	GetUserByUsername(string, interface{}) error
	GetUserByID(any, interface{}) error
	GetUserByEmail(string, interface{}) error
	VerifyAccount(string, sql.NullTime, interface{}) error
	UpdatePswdHash(string, any, interface{}) error
	UpdateVerHash(string, time.Time, any, interface{}) error
}
