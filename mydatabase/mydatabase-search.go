package mydatabase

// FindUserByEmail finding user by email
func FindUserByEmail(email string) (User, bool) {
	// return FindUserByField("email", email)
	// }

	// FindUserByField finds user by any one field
	// func FindUserByStringField(field string, value interface{}) (User, bool) {

	rows, err := GetDb().Query("SELECT * FROM users WHERE email=?", email)
	if err != nil {
		panic("error in sql select: " + err.Error())
	}

	var u User
	for rows.Next() {
		if err := rows.Scan(&u.ID, &u.Password, &u.Email, &u.IsEmailConfirmed,
			&u.ConfirmSecret, &u.Fio, &u.RoleID); err != nil {
			panic("Scan error:" + err.Error())
		}
		return u, true
	}
	return u, false
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

var cachedRole Role
