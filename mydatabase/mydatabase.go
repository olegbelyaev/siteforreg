package mydatabase

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/gocraft/dbr"
	ifErr "github.com/olegbelyaev/siteforreg/errorwrapper"
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

// Lecture - лекция
type Lecture struct {
	ID          int    `form:"id"`
	LocationID  int    `form:"location_id" binding:"required"`
	When        string `form:"when" binding:"required"`
	GroupName   string `form:"group_name" binding:"required"`
	MaxSeets    int    `form:"max_seets" binding:"required"`
	Name        string `form:"name" binding:"required"`
	Description string `form:"description" binding:"required"`
}

// GetDb - возвращает пул соединений с БД
// возможно выпилим (используй тогда GetDBRSession)
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

// GetConn - возвращает соединение с бд пакета sql
// Помни, что нужно делать defer conn.Close()
// Возможно, выпилим и будем юзать только GetDBRSession
func GetConn() *sql.Conn {
	conn, err := GetDb().Conn(Ctx)
	ifErr.Panic("connection error", err)
	return conn
}

// DBRConn - соединение с БД получаемое через GetDBR()
var DBRConn *dbr.Connection

// GetDBRConn - альтернатива getDB, использующая пакет dbr
// см. GetDBRSession
func GetDBRConn(log dbr.EventReceiver) *dbr.Connection {
	if DBRConn == nil {
		pass := os.Getenv("MYSQL_SECRET")
		if len(pass) == 0 {
			panic("MYSQL_SECRET is EMPTY! (set MYSQL_SECRET env var and run me again)")
		}
		var err error
		DBRConn, err = dbr.Open("mysql", fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/siteforeg", pass), nil)
		ifErr.Panic("DBR connection error", err)
		DBRConn.SetMaxOpenConns(10)
	}
	return DBRConn
}

// GetDBRSession - returns GetDBRConn(log).NewSession(log)
// Предпочтительный способ для получения сессии работы с бд и работы с пакетом dbr
func GetDBRSession(log dbr.EventReceiver) *dbr.Session {
	return GetDBRConn(nil).NewSession(nil)
}
