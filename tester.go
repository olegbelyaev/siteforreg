package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/siteforreg/mydatabase"
)

func usernameForm(c *gin.Context) {
	c.HTML(http.StatusOK, "template1.html", gin.H{})
}

func saveun(c *gin.Context) {
	c.HTML(http.StatusOK, "template1.html", gin.H{})
}

func newlocation(c *gin.Context) {
	var store = sessions.NewCookieStore([]byte("supersecret"))
	session, err := store.Get(c.Request, "session-name")
	if err != nil {
		panic(err.Error())
	}
	lastcname, _ := session.Values["lastlocname"]

	c.HTML(http.StatusOK, "tmp_valid_locations.html", gin.H{
		"info_msg": lastcname,
	})
}

func inslocation(c *gin.Context) {
	var l mydatabase.Location
	if err := c.ShouldBind(&l); err != nil {
		c.HTML(http.StatusOK, "templateAddLocation.html", gin.H{
			"error_msg": err.Error(),
		})
	}

	var store = sessions.NewCookieStore([]byte("supersecret"))
	session, err := store.Get(c.Request, "session-name")
	if err != nil {
		panic(err.Error())
	}
	session.Values["lastlocname"] = l.Name
	session.Save(c.Request, c.Writer)

	mydatabase.AddLocation(l)
	c.HTML(http.StatusOK, "templ_locations.html", gin.H{})

}

func startreg(c *gin.Context) {
	c.HTML(http.StatusOK, "registration.html", gin.H{})

}

func GenerateConfirmSecret() string {
	return "abraka"
}

func sendletter(confSecret string) {}

func endreg(c *gin.Context) {
	secret := GenerateConfirmSecret()
	email := c.PostForm("email")
	user := mydatabase.User{
		Email:         email,
		Password:      c.PostForm("password"),
		ConfirmSecret: secret,
		Fio:           c.PostForm("fio"),
		RoleID:        4,
	}
	_, ok := mydatabase.FindUserByEmail(email)
	if ok {
		c.HTML(http.StatusOK, "email_exists.html", gin.H{})
	} else {
		mydatabase.AddUser(user)
		sendletter(secret)
		c.HTML(http.StatusOK, "registration_end.html", gin.H{})
	}

}

func mainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "tmp_main.html", gin.H{})
}

func loginEnd(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	user, ok := mydatabase.FindUserByEmail(email)
	if !ok {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"error_msg": "This email not exists.",
		})
	} else {
		if password != user.Password {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"error_msg": "Password incorret.",
			})
		} else if !user.IsEmailConfirmed {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"error_msg": "You did't activate your account. Check out your email.",
			})
		} else {
			c.HTML(http.StatusOK, "tmp_main.html", gin.H{})
		}
	}
}

func main() {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	//router.LoagdHTMLFiles("templates/template1.html", "templates/template2.html")

	router.GET("/main", mainPage)

	router.GET("/form1", usernameForm)
	router.POST("/form1", saveun)

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"error_msg": "All is OK",
		})
	})
	router.POST("/login/end", loginEnd)

	router.GET("/registration/start", startreg)

	router.POST("/registration/end", endreg)
	locations := router.Group("/locations")
	{
		locations.GET("/new", newlocation)
		locations.POST("/insert", inslocation)
	}

	router.Run(":8180")
	sql.Drivers()
}
