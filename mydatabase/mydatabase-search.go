package mydatabase

// FindUserByEmail finding user
func FindUserByEmail(email string) (User, bool) {

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

var cachedUser User

// FindUserByEmailCached - ищет по email в БД и возвращает User, кэшируя значение
// для ускорения повторного получения
func FindUserByEmailCached(email string) (User, bool) {
	if cachedUser.ID == 0 {
		var ok bool
		cachedUser, ok = FindUserByEmail(email)
		if !ok {
			return cachedUser, false
		}
	}
	return cachedUser, true
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

// FindRoleByIDCached - ищет по ID Role в БД и возвращает, кэшируя значение
// для ускорения повторного получения
func FindRoleByIDCached(id int) (Role, bool) {
	if cachedRole.ID == 0 {
		var ok bool
		cachedRole, ok = FindRoleByID(id)
		if !ok {
			return cachedRole, false
		}
	}
	return cachedRole, true
}
