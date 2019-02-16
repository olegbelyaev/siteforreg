package mydatabase

import (
	"database/sql"
	"log"

	"github.com/gocraft/dbr"
	ifErr "github.com/olegbelyaev/siteforreg/errorwrapper"
)

// FindUserByEmail finding user by email
func FindUserByEmail(email string) (User, bool) {
	return FindUserByField("email", email)
}

// FindUserByField - finds user by any one field
func FindUserByField(field string, value interface{}) (User, bool) {
	var u User
	_, err := GetDBRSession(nil).Select("*").From("users").Where(field+"=?", value).Load(&u)
	ifErr.Log("Error while find user", err)
	return u, u.ID != 0 && err == nil
}

// Ticket - bind user and lecture, which hi will listen
type Ticket struct {
	ID        int
	UserID    int
	LectureID int
}

// UserLectureTicket - incasulates user, his ticket, lecture and location
type UserLectureTicket struct {
	User     User
	Lecture  Lecture
	Location Location
	Ticket   Ticket
}

// FindUserLectionTicketsByField - find tickets by userID and any one field
func FindUserLectionTicketsByField(userID int, field string, value interface{}) (lectures []UserLectureTicket) {
	var u User
	var le Lecture
	var lo Location
	var t Ticket

	// Минимальное условие - по юзеру:
	cond := dbr.Eq("tickets.user_id", userID)

	// если указано (непусто) field, значит - дополнительное по AND условие в WHERE:
	if len(field) > 0 {
		cond = dbr.And(cond, dbr.Eq(field, value))
	}

	sql := GetDBRSession(nil).Select("*").From("users").
		LeftJoin("tickets", "users.id=tickets.user_id").
		LeftJoin("lectures", "tickets.lecture_id=lectures.id").
		LeftJoin("locations", "lectures.location_id=locations.id").
		Where(cond)

	rows, err := sql.Rows()

	ifErr.Panic("error while find lections", err)

	for rows.Next() {
		err := rows.Scan(&u.ID, &u.Password, &u.Email, &u.Fio, &u.RoleID,
			&t.ID, &t.UserID, &t.LectureID,
			&le.ID, &le.LocationID, &le.When, &le.GroupName, &le.MaxSeets, &le.Name, &le.Description,
			&lo.ID, &lo.Name, &lo.Address)
		ifErr.Panic("cant scan lectures", err)

		lectures = append(lectures, UserLectureTicket{
			User:     u,
			Lecture:  le,
			Location: lo,
			Ticket:   t})

	}
	return lectures
}

// FindRoleByID - поиск роли по ее ID
func FindRoleByID(id int) (Role, bool) {
	var r Role

	err := GetDBRSession(nil).Select("*").From("roles").Where("id=?", id).Limit(1).LoadOne(&r)

	if err == dbr.ErrNotFound {
		return r, false
	}
	return r, true
}

// FindLocationsByField - ищет площадки по одному из полей или все, если поле пустое
func FindLocationsByField(field string, value interface{}) (locations []Location) {

	var rows *sql.Rows
	var err error
	if len(field) > 0 {
		// если имя поля непустое
		rows, err = GetDb().Query("SELECT * FROM locations WHERE "+field+"=?", value)
		if err != nil {
			panic("error in sql select: " + err.Error())
		}

	} else {
		// если имя поля пустое
		rows, err = GetDb().Query("SELECT * FROM locations")
		if err != nil {
			panic("error in sql select: " + err.Error())
		}
	}

	var l Location
	for rows.Next() {
		if err := rows.Scan(&l.ID, &l.Name, &l.Address); err != nil {
			panic("Scan error:" + err.Error())
		}
		locations = append(locations, l)
	}
	return
}

// FindUsersByField - ищет пользователей (field=""выбирает все записи)
func FindUsersByField(field string, value interface{}) (users []User) {

	sql := GetDBRSession(nil).Select("*").From("users")
	if len(field) > 0 {
		sql = sql.Where(dbr.Eq(field, value))
	}
	_, err := sql.Load(&users)
	ifErr.Panic("Can't find users", err)
	return users
}

// FindLocOrgsByField - поиск в таблице locorg по значению поля (или всех если field="")
func FindLocOrgsByField(field string, value interface{}) (locorgs []LocOrg) {

	var rows *sql.Rows
	var err error
	if len(field) > 0 {
		// если имя поля непустое
		sqlQuery := `SELECT l.*,u.* 
		FROM locorg lo
		LEFT JOIN users u ON lo.organizer_id=u.id
		LEFT JOIN locations l on lo.location_id=l.id
		WHERE ` + field + "=?"

		rows, err = GetDb().Query(sqlQuery, value)
		// panic(fmt.Sprintf(sqlQuery, value))
		if err != nil {
			panic("error in sqlQuery select: " + err.Error())
		}

	} else {
		// если имя поля пустое
		rows, err = GetDb().Query(`SELECT l.*,u.* FROM locorg lo
			LEFT JOIN users u ON lo.organizer_id=u.id
			LEFT JOIN locations l on lo.location_id=l.id`)
		if err != nil {
			panic("error in sqlQuery select: " + err.Error())
		}
	}

	var lo LocOrg
	for rows.Next() {
		if err := rows.Scan(&lo.Location.ID, &lo.Location.Name, &lo.Location.Address,
			&lo.Organizer.ID, &lo.Organizer.Password, &lo.Organizer.Email,
			&lo.Organizer.Fio, &lo.Organizer.RoleID); err != nil {
			// если запро с с left join то могут быть
			// висячие locorgs указывающие на никакую прощадку
			log.Printf("Skipped bad incostistent locorg record")
			continue
			// panic("Scan error:" + err.Error())
		}
		locorgs = append(locorgs, lo)
	}
	return
}

// FindLecturesByField - ищет лекции по имени и значению поля. Если field=="" ищет все лекции.
func FindLecturesByField(field string, value interface{}) (lectures []Lecture) {

	var rows *sql.Rows
	var err error
	if len(field) > 0 {
		// если имя поля непустое
		rows, err = GetDb().Query("SELECT * FROM lectures WHERE "+field+"=?", value)
		ifErr.Panic("Can't sql select (lectures): ", err)

	} else {
		// если имя поля пустое
		rows, err = GetDb().Query("SELECT * FROM lectures")
		ifErr.Panic("Can't sql select (all lectures): ", err)
	}

	var l Lecture
	for rows.Next() {
		if err := rows.Scan(&l.ID, &l.LocationID, &l.When, &l.GroupName,
			&l.MaxSeets, &l.Name, &l.Description); err != nil {
			panic("Scan error:" + err.Error())
		}
		lectures = append(lectures, l)
	}
	return
}

// LectureLocation - type with Lecture and his Location inside
type LectureLocation struct {
	Lecture  Lecture
	Location Location
}

// FindLecturesLocationsByField - finds lestures with locations
func FindLecturesLocationsByField(field string, value interface{}) (lectures []LectureLocation) {

	var le Lecture
	var lo Location

	sql := GetDBRSession(nil).Select("*").From("lectures").
		LeftJoin("locations", "lectures.location_id=locations.id")

	if len(field) > 0 {
		sql = sql.Where("lecture."+field+"=?", value)
	}

	rows, err := sql.Rows()
	ifErr.Panic("cant finf lectures + locations", err)

	for rows.Next() {
		err := rows.Scan(&le.ID, &le.LocationID, &le.When, &le.GroupName, &le.MaxSeets, &le.Name, &le.Description,
			&lo.ID, &lo.Name, &lo.Address)
		ifErr.Panic("scan error", err)
		lectures = append(lectures, LectureLocation{Lecture: le, Location: lo})
	}
	return lectures
}
