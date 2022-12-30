package db_mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gokch/ornn/db"
)

func NewDsn(id, pw, addr, port, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true", id, pw, addr, port, dbName)
}

func New(addr, port, id, pw, dbName string) (*db.Conn, error) {
	db := &db.Conn{}
	err := db.Connect("mysql", NewDsn(id, pw, addr, port, dbName), dbName)
	if err != nil {
		return nil, err
	}
	return db, nil
}
