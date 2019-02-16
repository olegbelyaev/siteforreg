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

	_, err := GetDBRSession(nil).
		InsertInto("locations").
		Columns("name", "address").
		Values(l.Name, l.Address).
		Exec()

	ifErr.Panic("Can't insert location", err)
}

// AddUser -- add new user to mydatabase
func AddUser(u User) (int64, error) {
	res, err := GetDBRSession(nil).InsertInto("users").
		Pair("id", u.ID).
		Pair("password", u.Password).
		Pair("email", u.Email).
		Pair("fio", u.Fio).
		Pair("role_id", u.RoleID).
		Exec()

	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

// AddLocOrg - добавляет связ площадки и организатора в БД
func AddLocOrg(locationID int, organiserID int) {
	res, err := GetDBRSession(nil).InsertInto("locorg").
		Pair("location_id", locationID).
		Pair("organizer_id", organiserID).
		Exec()
	ifErr.Panic("can't insert into locorg", err)
	_, err = res.LastInsertId()
	ifErr.Panic("can't get last insert id after insert to locorg", err)
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
	res, err := GetDBRSession(nil).InsertInto("roles").
		Columns("id", "name", "lvl").
		Values(r.ID, r.Name, r.Lvl).
		Exec()
	ifErr.Panic("can't insert into roles", err)
	_, err = res.LastInsertId()
	ifErr.Panic("can't get last insert id after insert to roles", err)
}

// AddLecture - adds lecture to db
func AddLecture(l Lecture) {
	_, err := GetDBRSession(nil).InsertInto("lectures").
		Pair("location_id", l.LocationID).
		Pair("when", l.When).
		Pair("group_name", l.GroupName).
		Pair("max_seets", l.MaxSeets).
		Pair("name", l.Name).
		Pair("description", l.Description).
		Exec()
	ifErr.Panic("can't insert into lectures", err)
}

// BuyTicket - user buy ticket for lecture
func BuyTicket(userID int, lectureID int) (ok bool) {

	_, err := GetDBRSession(nil).InsertInto("tickets").
		Pair("user_id", userID).
		Pair("lecture_id", lectureID).
		Exec()

	ifErr.Panic("Can't get last insert id", err)
	return true
}
