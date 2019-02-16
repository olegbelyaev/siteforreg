package mydatabase

import ifErr "github.com/olegbelyaev/siteforreg/errorwrapper"

// UpdateUser - иизменяет данные пользователя в БД
func UpdateUser(u User) {
	_, err := GetDBRSession(nil).Update("users").
		Set("password", u.Password).
		Set("email", u.Email).
		Exec()

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
		Exec()
	ifErr.Panic("can't update lectures", err)
}
