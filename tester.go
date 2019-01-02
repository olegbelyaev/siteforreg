package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/siteforreg/database"
)

func usernameForm(c *gin.Context) {
	c.HTML(http.StatusOK, "template1.html", gin.H{})
}

func saveun(c *gin.Context) {
	c.HTML(http.StatusOK, "template1.html", gin.H{})
}

func newlocation(c *gin.Context) {
	c.HTML(http.StatusOK, "templateAddLocation.html", gin.H{})
}

func inslocation(c *gin.Context) {
	locname := c.PostForm("locate_name")
	locadr := c.PostForm("locate_address")
	l := database.Location{ID: 0, Name: locname, Address: locadr}
	database.AddLocation(l)
	c.HTML(http.StatusOK, "templ_locations.html", gin.H{})

}

func main() {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	//router.LoagdHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/form1", usernameForm)
	router.POST("/form1", saveun)

	locations := router.Group("/locations")
	{
		locations.GET("/new", newlocation)
		locations.POST("/insert", inslocation)
	}

	router.Run(":8180")
	sql.Drivers()
}
