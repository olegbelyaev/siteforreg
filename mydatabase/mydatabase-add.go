package mydatabase

import (
	"os"
	"strconv"

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
		Pair("roles", u.Roles).
		Exec()

	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

// AddLocOrg - добавляет связь площадки и организатора в БД
func AddLocOrg(locationID int, organiserID int) {
	_, err := GetDBRSession(nil).InsertInto("locorg").
		Pair("location_id", locationID).
		Pair("organizer_id", organiserID).
		Exec()
	ifErr.Panic("can't insert into locorg", err)
}

// AddUserRole - role: +/-1 -for listener, +/-2 - for organizer, +/-4 - for admin
func AddUserRole(userID int, role int) {
	roleStr := strconv.Itoa(role)
	sqlSet := ""
	if role > 0 {
		// добавление role к существующим ролям:
		sqlSet = "roles|" + roleStr
	}
	if role < 0 {
		// вычитание role из существующих, только если существующие роли больше вычитаемой:
		sqlSet = "if(0+roles>" + roleStr + ",roles-" + roleStr + ",'')"
	}
	if len(sqlSet) == 0 {
		return
	}

	GetDBRSession(nil).Exec("UPDATE users SET roles="+sqlSet+" WHERE id=?", userID)
}

// AddInitAdmin - добавляет админа, если его нет в БД
func AddInitAdmin() {
	uu := FindUsersByIntRole(4)
	if len(uu) == 0 {
		// нет админов, создаем:
		_, err := AddUser(User{
			Email:    "admin",
			Fio:      "admin-fio",
			Roles:    4, //"admin",
			Password: os.Getenv("ADMIN_SECRET"),
		})
		ifErr.Panic("can't create init admin", err)
	}
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
