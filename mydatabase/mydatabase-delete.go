package mydatabase

import "fmt"

// // UpdateUser - иизменяет данные пользователя в БД
// func UpdateUser(u User) (rowsAffected int64, err error) {
// 	conn := GetConn()
// 	defer conn.Close()
// 	result, err := conn.ExecContext(Ctx,
// 		`UPDATE users SET password=?, email=?,
// 		 fio=?, role_id=? WHERE id=?`,
// 		u.Password, u.Email, u.Fio, u.RoleID,
// 		u.ID)
// 	if err != nil {
// 		return
// 	}
// 	rowsAffected, err = result.RowsAffected()
// 	return
// }

// DeleteLocation - удаляет площадку
func DeleteLocation(ID int) error {
	if ID <= 0 {
		return fmt.Errorf("ID must be >0")
	}
	// todo: проверить отсутствие на ней организаторов, и если есть , удалить их из locorg
	_, err := GetDb().Query(`DELETE FROM locations WHERE id=?`, ID)
	return err
}
