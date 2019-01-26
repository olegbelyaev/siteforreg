package mydatabase

import "fmt"

// DeleteLocation - удаляет площадку
func DeleteLocation(ID int) error {
	if ID <= 0 {
		return fmt.Errorf("ID must be >0")
	}
	// todo: проверить отсутствие на ней организаторов, и если есть , удалить их из locorg
	_, err := GetDb().Query(`DELETE FROM locations WHERE id=?`, ID)
	return err
}
