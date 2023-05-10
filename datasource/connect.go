// Package datasource @author uangi 2023-05
package datasource

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/real-uangi/pu55y/character"
	"github.com/real-uangi/pu55y/config"
	"github.com/real-uangi/pu55y/plog"
	"time"
)

// database/sql 级别的DB,Stmt都是并发安全的

var db *sql.DB

func InitDataSource(conf *config.Datasource) {
	var err error = nil
	var cs string
	cs = character.AppendAll(
		"host=", conf.Host,
		" port=", conf.Port,
		" user=", conf.User,
		" password=", conf.Password,
		" dbname=", conf.Database,
		" sslmode=disable",
	)
	db, err = sql.Open("postgres", cs)
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
