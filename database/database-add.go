package database

// No - test func
func No() string {
	return "no"
}

// AddLocation -- add location to database
func AddLocation(l Location) {
	conn := GetConn()
	defer conn.Close()
	conn.ExecContext(Ctx,
		"INSERT INTO locations (name,address) values(?,?)",
		l.Name, l.Address)
}

/*
rows, err := conn.QueryContext(ctx, "SELECT * FROM locations")
	if err != nil {
		panic("error in sql select")
	}

	defer rows.Close()

	locations := make([]Location, 0)
	for rows.Next() {
		var l Location
		if err := rows.Scan(&l.id, &l.name, &l.address); err != nil {
			panic("Scan error:" + err.Error())
		}
		locations = append(locations, l)
	}

	c.HTML(http.StatusOK, "templ_locations.html", gin.H{
		"locations": locations,
	})
*/
