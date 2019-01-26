package app

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/olegbelyaev/siteforreg/mydatabase"
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

// ShowMainPage - показ главной страницы сайта
func ShowMainPage(c *gin.Context) {
	c.Set("LoggedUser", GetLoggedUserFromSession(c))
	c.HTML(http.StatusOK, "main.html", c.Keys)
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
