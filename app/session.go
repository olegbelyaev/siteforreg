package app

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olegbelyaev/siteforreg/mydatabase"
	"github.com/olegbelyaev/siteforreg/mysession"
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

// AddWarningFlash - добавить warning-флеш-сообщение пользователю
func AddWarningFlash(c *gin.Context, msg string) {
	sess, _ := mysession.GetSession(c)
	sess.AddFlash("WARNING:" + msg)
	sess.Save(c.Request, c.Writer)
}

// AddInfoFlash - добавить info-флеш-сообщение пользователю
func AddInfoFlash(c *gin.Context, msg string) {
	sess, _ := mysession.GetSession(c)
	sess.AddFlash("INFO:" + msg)
	sess.Save(c.Request, c.Writer)
}

// GetFlashes - возвращает список предупреждающих и иформационных флеш-сообщений (и удаляет их из сессии)
func GetFlashes(c *gin.Context) ([]string, []string) {
	sess, _ := mysession.GetSession(c)
	allFlashes := sess.Flashes()
	sess.Save(c.Request, c.Writer)
	var warningFlashes []string
	var infoFlashes []string
	for _, fl := range allFlashes {
		if strings.HasPrefix(fl.(string), "WARNING") {
			warningFlashes = append(warningFlashes, fl.(string))
			continue
		}
		infoFlashes = append(infoFlashes, fl.(string))
	}

	return warningFlashes, infoFlashes
}
