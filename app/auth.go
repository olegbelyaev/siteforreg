package app

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/mail"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/olegbelyaev/siteforreg/mydatabase"
	"github.com/olegbelyaev/siteforreg/myemail"
	"github.com/olegbelyaev/siteforreg/mysession"
)

// GenerateSecret - generates random password
func GenerateSecret() string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		digits
	length := 10
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	for i := 1; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	return string(buf)
}

// GotoLoginIfNotLogged - подключается к роутеру через Use() до выполнения кода, требующего аутентификации
// "middleware" в терминах gin
func GotoLoginIfNotLogged(c *gin.Context) {
	c.Set("LoggedUser", GetLoggedUserFromSession(c))
	u, _ := c.Get("LoggedUser")
	if !u.(LoggedUser).IsLogged {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
}

// GotoAccessDeniedIfNotAdmin - проверяет уровень юзера >=4
// иначе редиректит на ошибку
func GotoAccessDeniedIfNotAdmin(c *gin.Context) {
	u := GetLoggedUserFromSession(c)
	// только в одном случае все ОК:
	if u.IsLogged && u.IsRoleFound && u.Role.Lvl >= 4 {
		return
	}
	// иначе редирект
	// todo: установка s.Set с редиректом не работает,
	// можно будет использовать флеш-сообщения через сесии
	// c.Redirect(http.StatusTemporaryRedirect, "/")
	// остановить цепочку
	c.Abort()
	// перенаправить на главную:
	c.Set("warning_msg", "Недостаточно прав")
	ShowMainPage(c)
}

// LoggedUser - IsLogged field answers on question "Is user logged?"
type LoggedUser struct {
	IsLogged    bool
	User        mydatabase.User
	IsRoleFound bool
	Role        mydatabase.Role
}

// GetLoggedUserRoleLvl - возвращает уровень роли пользователя
// возвращает 0 если роьне найдена или пользователь не залогинен
// todo: сейчас не проверить, т.к. нету пользователей с ролями меньшими 4
func GetLoggedUserRoleLvl(lu LoggedUser) int {
	if !lu.IsLogged {
		return 0
	}
	if !lu.IsRoleFound {
		return 0
	}
	return lu.Role.Lvl
}

// RegistrationEnd - обработка формы регистрации
func RegistrationEnd(c *gin.Context) {
	user := mydatabase.User{
		Email:    c.PostForm("email"),
		Password: GenerateSecret(),
		Fio:      c.PostForm("fio"),
		RoleID:   4,
	}
	_, ok := mydatabase.FindUserByEmail(user.Email)
	if ok {
		// Пользователь уже существует ( todo: добавить сброс пароля)
		c.HTML(http.StatusOK, "email_exists.html", c.Keys)
	} else {
		// TODO придумать как защититься от многочисленной отправки
		// пользователем писем по разным адресам
		siteAddress := "http://localhost:8081"
		emailText := fmt.Sprintf(
			`Здравствуйте, %s !
Вы зарегистрировались на сайте %s .
Ваш пароль: %s
			`, user.Fio, siteAddress, user.Password)

		sendErr := myemail.SendMailWithDefaultParams(
			mail.Address{Name: user.Fio, Address: user.Email},
			fmt.Sprintf("Регистрация на %s", siteAddress),
			emailText,
		)
		if sendErr != nil {
			c.Set("warning_msg", sendErr.Error())
			ShowMainPage(c)
			log.Print("Sending Error:" + sendErr.Error())
		}

		_, err := mydatabase.AddUser(user)
		if err != nil {
			log.Printf("Can't add user %v: %s", user, err.Error())
		}
		c.HTML(http.StatusOK, "main.html", c.Keys)
	}
}

// LoginEnd - обработка данных формы логина
func LoginEnd(c *gin.Context) {
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
			mysession.SaveEmail(c, email)
			// c.HTML(http.StatusOK, "main.html", gin.H{})
			ShowMainPage(c)
		}
	}
}
