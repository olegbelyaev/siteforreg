package app

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/mail"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/olegbelyaev/siteforreg/mydatabase"
	"github.com/olegbelyaev/siteforreg/myemail"
	"github.com/olegbelyaev/siteforreg/mysession"
)

// ShowMainPage - показ главной страницы сайта
func ShowMainPage(c *gin.Context) {
	c.Set("LoggedUser", GetLoggedUserFromSession(c))
	c.HTML(http.StatusOK, "main.html", c.Keys)
}

// SetWarningMsg - установка сообщения, которое пробросится в шаблоны ключем "warning_msg"
func SetWarningMsg(c *gin.Context, msg string) {
	c.Set("warning_msg", msg)
}

// ShowUsers -список пользователей
func ShowUsers(c *gin.Context) {
	c.Set("users", mydatabase.FindUsersByField("", ""))
	c.HTML(http.StatusOK, "show_users.html", c.Keys)
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

// ShowLocorgs - список
func ShowLocorgs(c *gin.Context) {
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

// ShowMyLocOrgs - список площадок текущего пользователя-организатора
func ShowMyLocOrgs(c *gin.Context) {
	curUser := GetLoggedUserFromSession(c)
	c.Set("locorgs", mydatabase.FindLocOrgsByField("organizer_id", curUser.User.ID))
	c.HTML(http.StatusOK, "show_locorgs.html", c.Keys)
}

// ShowMyLecures - shows lections manage page
func ShowMyLectures(c *gin.Context) {
	locationIDStr := c.Query("location_id")
	c.Set("LocationID", locationIDStr)
	c.Set("lectures", mydatabase.FindLecturesByField("location_id", locationIDStr))
	c.HTML(http.StatusOK, "manage_lectures.html", c.Keys)
}

// AddLocOrg - добавить организатора на площадку
func AddLocOrg(c *gin.Context) {
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
	ShowLocorgs(c)

}

// DeleteLocorg - удалить организатора на площадке
func DeleteLocorg(c *gin.Context) {
	locationIDStr := c.Query("location_id")
	locationID, err := strconv.Atoi(locationIDStr)
	if err != nil {
		c.Set("warning_msg", err.Error())
		ShowLocorgs(c)
		return
	}
	userIDStr := c.Query("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.Set("warning_msg", err.Error())
		ShowLocorgs(c)
		return
	}
	err = mydatabase.DeleteLocorgs(locationID, userID)
	if err != nil {
		c.Set("warning_msg", err.Error())
		ShowLocorgs(c)
		return
	}
	ShowLocorgs(c)
}

// ShowLocations - список площадок
func ShowLocations(c *gin.Context) {
	userID := c.Query("user_id")
	if len(userID) > 0 {
		c.Set("locations", mydatabase.FindLocationsByField(userID, ""))
	}
	c.Set("locations", mydatabase.FindLocationsByField("", ""))
	c.HTML(http.StatusOK, "locations.html", c.Keys)
}

// InsertLocation - вставка новой площадки
func InsertLocation(c *gin.Context) {
	var l mydatabase.Location
	// получение данных из формы создания новой площадки:
	if err := c.ShouldBind(&l); err != nil {
		c.Set("warning_msg", err.Error())
		ShowLocations(c)
	}
	mydatabase.AddLocation(l)

	// что лучше
	// c.Redirect(http.StatusTemporaryRedirect, "/administrate/locations/")
	ShowLocations(c)
}

// DeleteLocation - удаление площадки
func DeleteLocation(c *gin.Context) {
	locationIDStr := c.Param("ID")
	if len(locationIDStr) == 0 {
		c.Set("warning_msg", "ошибка данных")
		c.Abort()
		ShowLocations(c)
		return
	}
	locationID, err := strconv.Atoi(locationIDStr)
	if err != nil {
		c.Set("warning_msg", "ошибка парсинга данных ")
		c.Abort()
		ShowLocations(c)
		return
	}

	// поиск организаторов на этой площадке
	foundLocorgs := mydatabase.FindLocOrgsByField("location_id", locationID)
	if len(foundLocorgs) > 0 {
		// если на ней есть организаторы
		// сообщение пользователю и редирект на список организаторов этой площадки:
		mysession.AddInfoFlash(c, "Чтобы удалить площадку, удалите всех организаторов с данной лощадки")
		c.Redirect(http.StatusTemporaryRedirect, "/administrate/locorgs/?location_id="+locationIDStr)
		return
	}

	// удалить площадку
	err = mydatabase.DeleteLocation(locationID)
	if err != nil {
		c.Set("warning_msg", err.Error())
	}
	ShowLocations(c)
}

// EditLocation - форма редактирования и сохранение площадки
func EditLocation(c *gin.Context) {
	// какую площадку редактируем/сохраняем:
	locIDStr := c.Param("ID")
	locID, err := strconv.Atoi(locIDStr)
	if err != nil {
		c.Set("warning_msg", err.Error())
		// c.Abort()
		ShowLocations(c)
		return
	}

	// найдем ее в бд:
	locations := mydatabase.FindLocationsByField("id", locID)
	if len(locations) == 0 {
		// c.Abort()
		c.Set("warning_msg", "Not found")
		ShowLocations(c)
		return
	}

	// первая запись из найденных
	dbLocation := locations[0]

	// показываем форму редактирования, передаем туда данные из бд
	c.Set("Location", dbLocation)
	c.HTML(http.StatusOK, "edit_location.html", c.Keys)
	return
}

// SaveLocation - сохранение площадки
func SaveLocation(c *gin.Context) {

	// сохранение формы
	var formLocation mydatabase.Location
	if err := c.ShouldBind(&formLocation); err != nil {
		c.Set("warning_msg", err.Error())
		// c.Abort()
		ShowLocations(c)
		return
	}

	// ошибок нет, сохраняем данные в бд
	_, err := mydatabase.UpdateLocation(formLocation)
	if err != nil {
		c.Set("warning_msg", err.Error())
		// c.Abort()
	}
	ShowLocations(c)

}

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
