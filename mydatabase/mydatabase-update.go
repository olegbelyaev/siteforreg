package mydatabase

// UpdateUser - иизменяет данные пользователя в БД
func UpdateUser(u User) (rowsAffected int64, err error) {
	conn := GetConn()
	defer conn.Close()
	result, err := conn.ExecContext(Ctx,
		`UPDATE users SET password=?, email=?, 
		 fio=?, role_id=? WHERE id=?`,
		u.Password, u.Email, u.Fio, u.RoleID,
		u.ID)
	if err != nil {
		return
	}
	rowsAffected, err = result.RowsAffected()
	return
}