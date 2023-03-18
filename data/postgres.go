package data

import (
	"database/sql"
	"log"
)

type postgres struct {
	DB *sql.DB
}

func (p *postgres) New(pool *sql.DB) {
	log.Print(" postgres pool is running")
	p.DB = pool
}
