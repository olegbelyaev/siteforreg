package mydatabase

import (
	"os"

	ifErr "github.com/olegbelyaev/siteforreg/errorwrapper"
)

// No - test func
func No() string {
	return "no"
}

// AddLocation -- add location to mydatabase
func AddLocation(l Location) {
	conn := GetConn()
	defer conn.Close()
	conn.ExecContext(Ctx,
		"INSERT INTO locations (name,address) values(?,?)",
		l.Name, l.Address)
}

// AddUser -- add new user to mydatabase
func AddUser(u User) (int64, error) {
	conn := GetConn()
	defer conn.Close()
	result, err := conn.ExecContext(Ctx,
		`INSERT INTO users (id, password, email, fio, role_id)
		values(?,?,?,?,?)`,
		u.ID, u.Password, u.Email, u.Fio, u.RoleID)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// AddLocOrg - добавляет связ площадки и организатора в БД
func AddLocOrg(locationID int, organiserID int) {
	conn := GetConn()
	defer conn.Close()
	conn.ExecContext(Ctx,
		`INSERT INTO locorg (location_id, organizer_id)
	VALUES (?,?)`,
		locationID, organiserID)
}

// AddInitAdmin - добавляет админа, если его нет в БД
func AddInitAdmin() {
	_, ok := FindUserByField("role_id", 1)
	if !ok {
		_, err := AddUser(User{
			Email:    "admin",
			Fio:      "admin-fio",
			RoleID:   1,
			Password: os.Getenv("ADMIN_SECRET"),
		})
		ifErr.Panic("can't create init admin", err)
	}

}

// AddInitRoles - добавляет начальную роль суперадмина
// INSERT INTO roles (id,name,lvl) VALUES (1,"root", 4);
func AddInitRoles() {
	// admin:
	_, alreadyExists := FindRoleByID(1)
	if !alreadyExists {
		AddRole(Role{
			ID:   1,
			Name: "Администратор",
			Lvl:  "admin",
		})
	}
	// listener:
	_, alreadyExists = FindRoleByID(4)
	if !alreadyExists {
		AddRole(Role{
			ID:   4,
			Name: "Слушатель",
			Lvl:  "listener",
		})
	}
}

// AddRole - добавляет роль в БД
func AddRole(r Role) {
	conn := GetConn()
	defer conn.Close()
	_, err := conn.ExecContext(Ctx,
		"INSERT INTO roles (id, name, lvl) "+
			"VALUES (?,?,?)",
		r.ID, r.Name, r.Lvl)

	ifErr.Panic("Can't insert new lecture", err)
}

// AddLecture - adds lecture to db
func AddLecture(l Lecture) {
	conn := GetConn()
	defer conn.Close()
	_, err := conn.ExecContext(Ctx,
		"INSERT INTO lectures (location_id, `when`, group_name, max_seets, name, description) "+
			"VALUES (?,?,?,?,?,?)",
		l.LocationID, l.When, l.GroupName, l.MaxSeets, l.Name, l.Description)

	ifErr.Panic("Can't insert new lecture", err)
}

// SaveLecture - saves lecture to db
func SaveLecture(l Lecture) {
	conn := GetConn()
	defer conn.Close()
	_, err := conn.ExecContext(Ctx,
		"UPDATE lectures SET location_id=?, `when`=?, group_name=?, max_seets=?, name=?, description=? "+
			"WHERE id=?",
		l.LocationID, l.When, l.GroupName, l.MaxSeets, l.Name, l.Description, l.ID)

	ifErr.Panic("Can't save lecture", err)
}

// DeleteLecture - deletes lecture from db
func DeleteLecture(lectureID interface{}) {
	_, err := GetConn().ExecContext(Ctx,
		"DELETE FROM lectures WHERE id=?", lectureID)
	ifErr.Panic("Can't delete lecture ", err)
}

// BuyTicket - user buy ticket for lecture
func BuyTicket(userID int, lectureID int) (ok bool) {
	res, err := GetDBRSession(nil).Exec("INSERT INTO tickets (user_id, lecture_id) VALUES (?,?)", userID, lectureID)
	ifErr.Panic("can't buy ticket", err)
	_, err = res.LastInsertId()
	return err == nil
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
