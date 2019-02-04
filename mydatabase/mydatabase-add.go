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
		AddUser(User{
			Email:    "admin",
			Fio:      "admin-fio",
			RoleID:   1,
			Password: os.Getenv("ADMIN_SECRET"),
		})
	}

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
