package mydatabase

import "os"

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
func AddUser(u User) {
	conn := GetConn()
	defer conn.Close()
	conn.ExecContext(Ctx,
		`INSERT INTO users (id, password, email, is_email_confirmed, confirm_secret, fio, role_id)
		values(?,?,?,?,?,?,?)`,
		u.ID, u.Password, u.Email, u.IsEmailConfirmed, u.ConfirmSecret, u.Fio, u.RoleID)
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
			Email:            "admin",
			Fio:              "admin-fio",
			IsEmailConfirmed: true,
			RoleID:           1,
			Password:         os.Getenv("ADMIN_SECRET"),
		})
	}

}
