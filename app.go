package main

import (
	"encoding/gob"
	"github.com/ahmedkhaeld/booking/data"
	"github.com/ahmedkhaeld/booking/handlers"
	"github.com/ahmedkhaeld/booking/middleware"
	"github.com/ahmedkhaeld/jazz"
	"log"
	"os"
)

type application struct {
	*jazz.Jazz
	Handlers   *handlers.Handlers
	Models     data.Models
	Middleware *middleware.Middleware
}

func run() *application {
	gob.Register(data.Reservation{})
	gob.Register(data.User{})
	gob.Register(data.Restriction{})
	gob.Register(data.Room{})

	//get the root path of the application
	rootPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	//define a jazz instance
	j := &jazz.Jazz{}

	//init the jazz
	err = j.New(rootPath)
	if err != nil {
		j.ErrorLog.Fatal(err)
	}

	m := &middleware.Middleware{
		Jazz: j,
	}

	h := &handlers.Handlers{
		Jazz: j,
	}
	app := &application{
		Jazz:       j,
		Handlers:   h,
		Middleware: m,
	}
	app.Models = data.New(app.Jazz.DB.SqlPool)
	app.Middleware.Models = app.Models
	//add new routes with default ones
	app.Jazz.Routes = app.routes()
	return app
}
