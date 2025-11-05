package vrc_world_matching

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

type User struct {
	ID string `db:"first_name"`
}

func init() {
	var err error
	db, err = sqlx.Connect("mysql", "root:@(localhost:4000)/test")
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
}
