package data

import (
	"database/sql"
	"os"
)

var DB *sql.DB
var psql postgres

type Models struct {
	//any models inserted here (and in the New func)
	//are easily accessible throughout the entire application
	Rooms        Room
	Users        User
	Reservations Reservation
	Restrictions Restriction
}

func New(databasePool *sql.DB) Models {
	DB = databasePool

	switch os.Getenv("DATABASE_TYPE") {
	case "postgres", "postgresql":
		psql.New(databasePool)
	default:
		//do nothing
	}

	return Models{
		Rooms:        Room{},
		Users:        User{},
		Reservations: Reservation{},
		Restrictions: Restriction{},
	}
}
