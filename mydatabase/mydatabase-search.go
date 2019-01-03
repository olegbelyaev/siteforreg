package mydatabase

// FindUserByEmail finding user
func FindUserByEmail(email string) (User, bool) {

	rows, err := GetDb().Query("SELECT * FROM users WHERE email=?", email)
	if err != nil {
		panic("error in sql select")
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
