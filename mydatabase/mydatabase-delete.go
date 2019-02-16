package mydatabase

import (
	"fmt"

	"github.com/gocraft/dbr"
	ifErr "github.com/olegbelyaev/siteforreg/errorwrapper"
)

// DeleteLocation - удаляет площадку
func DeleteLocation(ID int) error {
	if ID <= 0 {
		return fmt.Errorf("ID must be >0")
	}
	// todo: проверить отсутствие на ней организаторов, и если есть , удалить их из locorg
	_, err := GetDBRSession(nil).DeleteFrom("locations").Where("id=?", ID).Exec()
	return err
}

// DeleteLocorgs - удалить одну или несколько записей locorg
func DeleteLocorgs(locationID int, userID int) error {
	// _, err := GetDb().Query(`DELETE FROM locorg WHERE location_id=? AND organizer_id=?`, locationID, userID)
	_, err := GetDBRSession(nil).DeleteFrom("locorg").
		Where(
			dbr.And(
				dbr.Eq("location_id", locationID),
				dbr.Eq("organizer_id", userID))).
		Exec()

	return err
}

// ReleaseTicket - удалить билет из БД
func ReleaseTicket(ticketID int, userID int) bool {
	// удалить из БД ticketID для userID
	res, err := GetDBRSession(nil).DeleteFrom("tickets").
		Where("id=? AND user_id=?", ticketID, userID).Limit(1).Exec()
	ifErr.Panic("Can't exec delete ticket", err)
	affected, err := res.RowsAffected()
	ifErr.Panic("Can't get rows affected after deleting ticket", err)
	if affected == 0 {
		return false
	}
	return true
}

// DeleteLecture - deletes lecture from db
func DeleteLecture(lectureID interface{}) {
	_, err := GetDBRSession(nil).DeleteFrom("lectures").
		Where(dbr.Eq("id", lectureID)).
		Exec()
	ifErr.Panic("can't delete from lectures", err)
}
