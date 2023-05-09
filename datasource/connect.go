package datasource

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/real-uangi/pu55y/character"
	"github.com/real-uangi/pu55y/plog"
	"time"
)

// database/sql 级别的DB,Stmt都是并发安全的

var db *sql.DB

func InitDataSource(host string, port string, usr string, pw string, dbn string) {
	var err error = nil
	var config string
	config = character.AppendAll("host=", host, " port=", port, " user=", usr, " password=", pw, " dbname=", dbn, " sslmode=disable")
	db, err = sql.Open(
		"postgres",
		config,
	)
	if err != nil {
		plog.Error(err.Error())
	} else {
		db.SetConnMaxLifetime(time.Hour)
		db.SetConnMaxIdleTime(5 * time.Minute)
		db.SetMaxIdleConns(2)
		db.SetMaxOpenConns(8)
	}
}

func GetDataSource() *sql.DB {
	return db
}
