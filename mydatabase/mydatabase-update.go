package mydatabase

import (
	ifErr "github.com/olegbelyaev/siteforreg/errorwrapper"
)

// UpdateUser - изменяет данные пользователя в БД
func UpdateUser(u User) {

	stmt := GetDBRSession(nil).Update("users").
		Set("password", u.Password).
		Set("reset_key", u.ResetKey).
		Set("email", u.Email).
		Set("fio", u.Fio).
		Set("roles", u.Roles).
		Where("id=?", u.ID)

	// panic(fmt.Sprintf("%v", stmt.Query))

	_, err := stmt.Exec()

	ifErr.Panic("can't update user", err)
}

// UpdateLocation - иизменяет данные площадки в БД
func UpdateLocation(l Location) {
	_, err := GetDBRSession(nil).Update("locations").
		Set("name", l.Name).
		Set("address", l.Address).
		Where("id=?", l.ID).
		Exec()
	ifErr.Panic("can't update locations", err)
}

// SaveLecture - saves lecture to db
func SaveLecture(l Lecture) {
	_, err := GetDBRSession(nil).Update("lectures").
		Set("location_id", l.LocationID).
		Set("when", l.When).
		Set("group_name", l.GroupName).
		Set("max_seets", l.MaxSeets).
		Set("name", l.Name).
		Set("description", l.Description).
		Where("id=?", l.ID).
		Exec()
	ifErr.Panic("can't update lectures", err)
}
