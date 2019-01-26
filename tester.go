package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strconv"

	"github.com/olegbelyaev/siteforreg/app"
	"github.com/olegbelyaev/siteforreg/myemail"

	"github.com/olegbelyaev/siteforreg/mysession"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/olegbelyaev/siteforreg/mydatabase"
)

func usernameForm(c *gin.Context) {
	c.HTML(http.StatusOK, "template1.html", c.Keys)
}

func saveun(c *gin.Context) {
	c.HTML(http.StatusOK, "template1.html", c.Keys)
}

func newlocation(c *gin.Context) {
	session, _ := mysession.GetSession(c)
	lastcname, _ := session.Values["lastlocname"]
	c.Set("info_msg", lastcname)
	c.HTML(http.StatusOK, "tmp_valid_locations.html", c.Keys)
}

func showLocations(c *gin.Context) {
	userID := c.Query("user_id")
	if len(userID) > 0 {
		c.Set("locations", mydatabase.FindLocationsByField(userID, ""))
	}
	c.Set("locations", mydatabase.FindLocationsByField("", ""))
	c.HTML(http.StatusOK, "locations.html", c.Keys)
}

func showUsers(c *gin.Context) {
	c.Set("users", mydatabase.FindUsersByField("", ""))
	c.HTML(http.StatusOK, "show_users.html", c.Keys)
}

func showLocorgs(c *gin.Context) {
	locationID := c.Query("location_id")
	userID := c.Query("user_id")
	var locorgs []mydatabase.LocOrg
	if len(locationID) > 0 {
		locorgs = mydatabase.FindLocOrgsByField("location_id", locationID)
	} else if len(userID) > 0 {
		locorgs = mydatabase.FindLocOrgsByField("organizer_id", userID)
	} else {
		locorgs = mydatabase.FindLocOrgsByField("", "")
	}
	c.Set("locorgs", locorgs)
	c.HTML(http.StatusOK, "show_locorgs.html", c.Keys)
}

func inslocation(c *gin.Context) {
	var l mydatabase.Location
	// получение данных из формы создания новой площадки:
	if err := c.ShouldBind(&l); err != nil {
		c.Set("warning_msg", err.Error())
		c.HTML(http.StatusOK, "templateAddLocation.html", c.Keys)
	}

	mydatabase.AddLocation(l)
	c.Redirect(http.StatusTemporaryRedirect, "/locations/")
}

func dellocation(c *gin.Context) {
	locationIDStr := c.Param("ID")
	if len(locationIDStr) == 0 {
		c.Set("warning_msg", "ошибка данных")
		c.Abort()
		showLocations(c)
		return
	}
	locationID, err := strconv.Atoi(locationIDStr)
	if err != nil {
		c.Set("warning_msg", "ошибка парсинга данных ")
		c.Abort()
		showLocations(c)
		return
	}

	err = mydatabase.DeleteLocation(locationID)
	if err != nil {
		c.Set("warning_msg", err.Error())
	}
	showLocations(c)
}

// вызывает форму регистрации на сайте
func startreg(c *gin.Context) {
	c.HTML(http.StatusOK, "registration.html", c.Keys)
}

// когда пользователь заполнил форму регистрации нового юзера на сайт
func endreg(c *gin.Context) {
	email := c.PostForm("email")
	secret := app.GenerateSecret()
	user := mydatabase.User{
		Email:    email,
		Password: secret,
		Fio:      c.PostForm("fio"),
		RoleID:   4,
	}
	_, ok := mydatabase.FindUserByEmail(email)
	if ok {
		// todo: Пользователь уже существует, здесь должна быть ссылка или форма на сброс пароля
		c.HTML(http.StatusOK, "email_exists.html", c.Keys)
	} else {
		// TODO придумать как защититься от многочисленной отправки
		// пользователем писем по разным адресам
		siteAddress := "http://localhost:8081"
		emailText := fmt.Sprintf(
			`Здравствуйте, %s !
Вы зарегистрировались на сайте %s .
Ваш пароль: %s
			`, user.Fio, siteAddress, secret)

		sendErr := myemail.SendMailWithDefaultParams(
			mail.Address{Name: user.Fio, Address: user.Email},
			fmt.Sprintf("Регистрация на %s", siteAddress),
			emailText,
		)
		if sendErr != nil {
			// TODO здесь вместо паники, выводить сообщение пользователю
			// на страницу
			panic("Sending Error:" + sendErr.Error())
		}

		// panic("----------------OK-----------------")

		_, err := mydatabase.AddUser(user)
		if err != nil {
			log.Printf("Can't add user %v: %s", user, err.Error())
		}
		c.HTML(http.StatusOK, "main.html", c.Keys)
	}
}

// пользователь заполнил и отправил форму входа юзера на сайт
func loginEnd(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	user, ok := mydatabase.FindUserByEmail(email)
	if !ok {
		c.Set("warning_msg", "This email not exists.")
		c.HTML(http.StatusOK, "login.html", c.Keys)
	} else {
		if password != user.Password {
			// юзер найден но пароль не совпадает:
			c.Set("warning_msg", "Password incorret.")
			c.HTML(http.StatusOK, "login.html", c.Keys)
		} else {
			// юзер найден и емаил подтвержден:
			app.SaveEmailToSession(c, email)
			// c.HTML(http.StatusOK, "main.html", gin.H{})
			app.ShowMainPage(c)
		}
	}
}

// разлогинить пользователя
func logout(c *gin.Context) {
	// для разлогина сохраним емаил, по которому пользователь не найдется
	app.SaveEmailToSession(c, "")
	app.ShowMainPage(c)
}

func addLocOrg(c *gin.Context) {
	locID := c.PostForm("location_id")
	if len(locID) > 0 {
		c.Set("location_id", locID)
	}

	userID := c.PostForm("user_id")
	if len(userID) > 0 {
		c.Set("user_id", userID)
	}

	if len(locID) == 0 {
		c.Set("locations", mydatabase.FindLocationsByField("", ""))
		c.HTML(http.StatusOK, "select_location.html", c.Keys)
		return
	}
	if len(userID) == 0 {
		c.Set("users", mydatabase.FindUsersByField("", ""))
		// panic(fmt.Sprintf("------------------>%s", c.Keys))
		c.HTML(http.StatusOK, "select_user.html", c.Keys)
		return
	}
	locIDint, err := strconv.Atoi(locID)
	if err != nil {
		panic("Can't parse as int:" + locID)
	}
	userIDint, err := strconv.Atoi(userID)
	if err != nil {
		panic("Can't parse as int:" + userID)
	}
	mydatabase.AddLocOrg(locIDint, userIDint)

}

func main() {

	mydatabase.AddInitAdmin()

	myemail.SetParams(
		"", "sivsite@yandex.ru", os.Getenv("EMAIL_SECRET"),
		"smtp.yandex.ru", "465",
		mail.Address{Name: "sitename", Address: "sivsite@yandex.ru"},
		true,
	)

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	//router.LoagdHTMLFiles("templates/template1.html", "templates/template2.html")
	// router.Static

	router.Use(func(c *gin.Context) {
		c.Set("html_title", "Siteforrreg")
		c.Set("LoggedUser", app.GetLoggedUserFromSession(c))
	})

	router.GET("/", app.ShowMainPage)

	router.GET("/form1", usernameForm)
	router.POST("/form1", saveun)

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", c.Keys)
	})
	router.POST("/login/end", loginEnd)

	router.GET("/logout", logout)

	router.GET("/registration/start", startreg)

	router.POST("/registration/end", endreg)

	locations := router.Group("/locations")
	{
		// логина не требует:
		locations.Any("/", showLocations)

		// ниже этого будет требовать залогиниться:
		locations.Use(app.GotoLoginIfNotLogged)
		// требование быть админом:
		locations.Use(app.GotoAccessDeniedIfNotAdmin)

		locations.GET("/new", newlocation)
		locations.POST("/insert", inslocation)
		// td: как защититься от запросов не с этого сайта?
		locations.Any("/delete/:ID", dellocation)
	}

	users := router.Group("/users")
	{
		users.Use(app.GotoLoginIfNotLogged)
		users.Use(app.GotoAccessDeniedIfNotAdmin)
		users.GET("/", showUsers)
	}

	router.GET("/locorgs", showLocorgs)
	router.Any("/add_locorg", addLocOrg)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		panic("PORT env not defined!")
	}

	router.Run(fmt.Sprintf(":%s", port))
	sql.Drivers()
}
