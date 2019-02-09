package mydatabase

import (
	"database/sql"
	"log"

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
	return u, err == nil
}

// FindRoleByID - поиск роли по ее ID
func FindRoleByID(id int) (Role, bool) {
	rows, err := GetDb().Query("SELECT * FROM roles WHERE id=?", id)
	if err != nil {
		panic("error in select: " + err.Error())
	}
	var r Role
	for rows.Next() {
		if err := rows.Scan(&r.ID, &r.Name, &r.Lvl); err != nil {
			panic("Scan error:" + err.Error())
		}
		return r, true
	}
	return r, false
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

	var rows *sql.Rows
	var err error
	if len(field) > 0 {
		// если имя поля непустое
		rows, err = GetDb().Query("SELECT * FROM users WHERE "+field+"=?", value)
		if err != nil {
			panic("error in sql select: " + err.Error())
		}

	} else {
		// если имя поля пустое
		rows, err = GetDb().Query("SELECT * FROM users")
		if err != nil {
			panic("error in sql select: " + err.Error())
		}
	}

	var u User
	for rows.Next() {
		if err := rows.Scan(&u.ID, &u.Password, &u.Email,
			&u.Fio, &u.RoleID); err != nil {
			panic("Scan error:" + err.Error())
		}
		users = append(users, u)
	}
	return
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
