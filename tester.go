package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/mail"
	"os"
	"strconv"

	"github.com/olegbelyaev/siteforreg/myemail"

	"github.com/olegbelyaev/siteforreg/mysession"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/olegbelyaev/siteforreg/mydatabase"
)

// SaveEmailToSession  - сохраняет email пользователя в сессию
func SaveEmailToSession(c *gin.Context, email string) {
	sess, _ := mysession.GetSession(c)
	sess.Values["email"] = email
	sess.Save(c.Request, c.Writer)
}

// GetUserFromSession - достает объект User из сессии
func GetUserFromSession(c *gin.Context) (mydatabase.User, bool) {
	var user mydatabase.User
	email, ok := mysession.GetStringValue(c, "email")
	if !ok || len(email) == 0 {
		return user, false
	}
	user, ok = mydatabase.FindUserByEmail(email)
	return user, ok
}

// LoggedUser - IsLogged field answers on question "Is user logged?"
type LoggedUser struct {
	IsLogged    bool
	User        mydatabase.User
	IsRoleFound bool
	Role        mydatabase.Role
}

// GetLoggedUserFromSession - получает из сессии email юзера.
// Создает структуру User поиском в БД по email. Если не найден - ставит IsLogged=false.
// Сохраняет юзера в поле User.
// По полю User.RoleID находит роль в БД. Если не найден ставит IsRoleFound=false.
// Сохраняет роль в поле Role.
func GetLoggedUserFromSession(c *gin.Context) LoggedUser {
	userFromSess, ok := GetUserFromSession(c)
	var lu = LoggedUser{
		IsLogged: ok,
		User:     userFromSess,
	}
	if lu.IsLogged {
		lu.Role, lu.IsRoleFound = mydatabase.FindRoleByID(lu.User.RoleID)
	}
	return lu
}

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
	if err := c.ShouldBind(&l); err != nil {
		c.Set("error_msg", err.Error())
		c.HTML(http.StatusOK, "templateAddLocation.html", c.Keys)
	}

	var store = sessions.NewCookieStore([]byte("supersecret"))
	session, err := store.Get(c.Request, "session-name")
	if err != nil {
		panic(err.Error())
	}
	session.Values["lastlocname"] = l.Name
	session.Save(c.Request, c.Writer)

	mydatabase.AddLocation(l)
	c.HTML(http.StatusOK, "templ_locations.html", c.Keys)
}

// вызывает форму регистрации на сайте
func startreg(c *gin.Context) {
	c.HTML(http.StatusOK, "registration.html", c.Keys)
}

// GenerateConfirmSecret - generates email comfirm secret parameter
func GenerateConfirmSecret() string {
	return "abraka"
}

// когда пользователь заполнил форму регистрации нового юзера на сайт
func endreg(c *gin.Context) {
	email := c.PostForm("email")
	secret := GenerateConfirmSecret()
	user := mydatabase.User{
		Email:         email,
		Password:      c.PostForm("password"),
		ConfirmSecret: secret,
		Fio:           c.PostForm("fio"),
		RoleID:        4,
	}
	_, ok := mydatabase.FindUserByEmail(email)
	if ok {
		// todo: здесь должна быть ссылка или форма на сброс пароля
		c.HTML(http.StatusOK, "email_exists.html", c.Keys)
	} else {
		// TODO придумать как защититься от многочисленной отправки
		// пользователем (одной кукой?) писем по разным адресам
		siteAddress := "http://localhost:8081"
		emailText := fmt.Sprintf(
			`Здравствуйте, %s !
			Вы зарегистрировались на сайте %s .
			Для подтверждения почты, перейдите по ссылке %s/confirmemail/%s`,
			user.Fio, siteAddress, siteAddress, secret)

		sendErr := myemail.SendMailWithDefaultParams(
			mail.Address{Name: user.Fio, Address: user.Email},
			"Регистрация", emailText)
		if sendErr != nil {
			// TODO здесь вместо паники, выводить сообщение пользователю
			// на страницу
			panic("Sending Error:" + sendErr.Error())
		}

		// panic("----------------OK-----------------")

		mydatabase.AddUser(user)
		c.HTML(http.StatusOK, "registration_end.html", c.Keys)
	}
}

func confirmemail(c *gin.Context) {
	confirmSecret := c.Param("secret")
	if len(confirmSecret) < 6 {
		c.Set("error_msg", "Код неправильный")
		c.HTML(http.StatusOK, "main.html", c.Keys)
		return
	}
	user, ok := mydatabase.FindUserByField("confirm_secret", confirmSecret)
	if !ok {
		c.Set("error_msg", "Код неправильный!")
		c.HTML(http.StatusOK, "main.html", c.Keys)
		return
	}
	user.IsEmailConfirmed = true
	user.ConfirmSecret = ""
	_, err := mydatabase.UpdateUser(user)
	if err != nil {
		// todo: вместо паники что-то придумать
		panic(err.Error())
	}
	c.Set("info_msg", "Email подтвержден.")
	c.HTML(http.StatusOK, "main.html", c.Keys)
}

// при переходе на главную страницу сайта
func showMainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "main.html", c.Keys)
}

// пользователь заполнил и отправил форму входа юзера на сайт
func loginEnd(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	user, ok := mydatabase.FindUserByEmail(email)
	if !ok {
		c.Set("error_msg", "This email not exists.")
		c.HTML(http.StatusOK, "login.html", c.Keys)
	} else {
		if password != user.Password {
			// юзер найден но пароль не совпадает:
			c.Set("error_msg", "Password incorret.")
			c.HTML(http.StatusOK, "login.html", c.Keys)
		} else if !user.IsEmailConfirmed {
			// юзер найден емаил не подтвержден:
			c.Set("error_msg", "You did't activate your account. Check out your email.")
			c.HTML(http.StatusOK, "login.html", c.Keys)
		} else {
			// юзер найден и емаил подтвержден:
			SaveEmailToSession(c, email)
			// c.HTML(http.StatusOK, "main.html", gin.H{})
			showMainPage(c)
		}
	}
}

// разлогинить пользователя
func logout(c *gin.Context) {
	// для разлогина сохраним емаил, по которому пользователь не найдется
	SaveEmailToSession(c, "")
	showMainPage(c)
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

	router.Use(func(c *gin.Context) {
		c.Set("html_title", "Siteforrreg")
		c.Set("LoggedUser", GetLoggedUserFromSession(c))
	})

	router.GET("/", showMainPage)

	router.GET("/form1", usernameForm)
	router.POST("/form1", saveun)

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"error_msg": "All is OK",
		})
	})
	router.POST("/login/end", loginEnd)

	router.GET("/logout", logout)

	router.GET("/registration/start", startreg)

	router.POST("/registration/end", endreg)

	router.GET("/confirmemail/:secret", confirmemail)

	locations := router.Group("/locations")
	{
		locations.GET("/", showLocations)
		locations.GET("/new", newlocation)
		locations.POST("/insert", inslocation)
	}

	router.GET("/users", showUsers)

	router.GET("/locorgs", showLocorgs)
	router.Any("/add_locorg", addLocOrg)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		panic("PORT env not defined!")
	}

	router.Run(fmt.Sprintf(":%s", port))
	sql.Drivers()
}
