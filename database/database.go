package database

import (
	"context"
	"database/sql"
)

//Db -- pull of connections
var Db *sql.DB

//Ctx -- we dont know
var Ctx = context.Background()

// Location - тип площадки
type Location struct {
	ID      int
	Name    string
	Address string
}

// GetDb - возвращает пул соединений с БД
func GetDb() *sql.DB {
	if Db == nil {
		db1, err := sql.Open("mysql", "root:11@tcp(127.0.0.1:3306)/siteforeg")
		if err != nil {
			panic("pool error:" + err.Error())
		}
		Db = db1
	}
	return Db
}

// GetConn - возвращает соединение с бд
// Помни, что нужно делать defer conn.Close()
func GetConn() *sql.Conn {
	conn, err := GetDb().Conn(Ctx)
	if err != nil {
		panic("connection error:" + err.Error())
	}
	return conn
}
