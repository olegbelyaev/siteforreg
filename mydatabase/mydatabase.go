package mydatabase

import (
	"context"
	"database/sql"
	"fmt"
	"os"
)

//Db -- pull of connections
var Db *sql.DB

//Ctx -- we dont know
var Ctx = context.Background()

// Location - тип площадки
type Location struct {
	ID      int
	Name    string `form:"locate_name" binding:"required"`
	Address string `form:"locate_address" binding:"required"`
}

// User - represent User
type User struct {
	ID       int
	Email    string
	Password string
	Fio      string
	RoleID   int
}

// LocOrg - represent Locations binded with Organizers
type LocOrg struct {
	Location  Location
	Organizer User
}

// Role - represent user role
type Role struct {
	ID   int
	Name string
	Lvl  int
}

// GetDb - возвращает пул соединений с БД
func GetDb() *sql.DB {
	if Db == nil {
		pass := os.Getenv("MYSQL_SECRET")
		if len(pass) == 0 {
			panic("MYSQL_SECRET is EMPTY! (set MYSQL_SECRET env var and run me again)")
		}

		db1, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/siteforeg", pass))
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
