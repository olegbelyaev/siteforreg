package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/mail"
	"os"

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
	c.HTML(http.StatusOK, "template1.html", gin.H{})
}

func saveun(c *gin.Context) {
	c.HTML(http.StatusOK, "template1.html", gin.H{})
}

func newlocation(c *gin.Context) {
	session, _ := mysession.GetSession(c)
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

// вызывает форму регистрации на сайте
func startreg(c *gin.Context) {
	c.HTML(http.StatusOK, "registration.html", gin.H{})

}

// GenerateConfirmSecret - generates email comfirm secret parameter
func GenerateConfirmSecret() string {
	return "abraka"
}

func sendletter(confSecret string) {}

// когда пользователь заполнил форму регистрации нового юзера на сайт
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
		// TODO придумать как защититься от многочисленной отправки
		// пользователем (одной кукой?) писем по разным адресам
		sendErr := myemail.SendMailWithDefaultParams(
			// подставить адрес и фио  пользователя из формы
			mail.Address{Name: "ФИО", Address: "d.l.belyaev@gmail.com"},
			"Регистрация",
			"Это текст сообщения ")
		if sendErr != nil {
			// TODO здесь вместо паники, выводить сообщение пользователю
			// на страницу
			panic("Sending Error:" + sendErr.Error())
		}

		// panic("----------------OK-----------------")

		c.HTML(http.StatusOK, "registration_end.html", gin.H{})
	}
}

// при переходе на главную страницу сайта
func showMainPage(c *gin.Context) {
	DefaultH["LoggedUser"] = GetLoggedUserFromSession(c)
	c.HTML(http.StatusOK, "tmp_main.html", DefaultH)
}

// пользователь заполнил и отправил форму входа юзера на сайт
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
			// юзер найден но пароль не совпадает:
			c.HTML(http.StatusOK, "login.html", gin.H{
				"error_msg": "Password incorret.",
			})
		} else if !user.IsEmailConfirmed {
			// юзер найден емаил не подтвержден:
			c.HTML(http.StatusOK, "login.html", gin.H{
				"error_msg": "You did't activate your account. Check out your email.",
			})
		} else {
			// юзер найден и емаил подтвержден:
			SaveEmailToSession(c, email)
			// c.HTML(http.StatusOK, "tmp_main.html", gin.H{})
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

// DefaultH - набор параметров по умолчанию для передачи в шаблоны
var DefaultH = make(map[string]interface{})

func main() {

	mydatabase.AddInitAdmin()

	myemail.SetParams(
		"", "sivsite@yandex.ru", os.Getenv("EMAIL_SECRET"),
		"smtp.yandex.ru", "465",
		mail.Address{Name: "sitename", Address: "sivsite@yandex.ru"},
	)

	DefaultH["aaa"] = "привет"
	// DefaultH["HasUserFromSessionLevelUpTo"] = HasUserFromSessionLevelUpTo

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	//router.LoagdHTMLFiles("templates/template1.html", "templates/template2.html")

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
	locations := router.Group("/locations")
	{
		locations.GET("/new", newlocation)
		locations.POST("/insert", inslocation)
	}

	port := os.Getenv("PORT")
	if len(port) == 0 {
		panic("PORT env not defined!")
	}

	router.Run(fmt.Sprintf(":%s", port))
	sql.Drivers()
}
