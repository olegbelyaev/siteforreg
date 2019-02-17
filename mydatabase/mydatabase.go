package mydatabase

import (
	"context"
	"fmt"
	"os"

	"github.com/gocraft/dbr"
	ifErr "github.com/olegbelyaev/siteforreg/errorwrapper"
)

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
	Roles    int
}

// UserRoleListener - слушательская составляющая поля roles таблицы users
const UserRoleListener = 1

// UserRoleOrganizer - организаторская составляющая поля roles таблицы users
const UserRoleOrganizer = 2

// UserRoleAdmin - админская составляющая поля roles таблицы users
const UserRoleAdmin = 4

// LocOrg - represent Locations binded with Organizers
type LocOrg struct {
	Location  Location
	Organizer User
}

// Role - represent user role
type Role struct {
	ID   int
	Name string
	Lvl  string
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
