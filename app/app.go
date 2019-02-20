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
	ifErr "github.com/olegbelyaev/siteforreg/errorwrapper"
	"github.com/olegbelyaev/siteforreg/mydatabase"
	"github.com/olegbelyaev/siteforreg/myemail"
	"github.com/olegbelyaev/siteforreg/mysession"
)

// ShowMainPage - показ главной страницы сайта
func ShowMainPage(c *gin.Context) {
	loggedUser := GetLoggedUserFromSession(c)
	log.Printf("%v", loggedUser)

	leclocs := mydatabase.FindLecturesLocationsByField("", "")
	c.Set("leclocs", leclocs)

	if loggedUser.IsLogged {
		c.Set("LoggedUser", loggedUser)

		c.HTML(http.StatusOK, "main_logged.html", c.Keys)
		return
	}

	// незалогированый юзер
	c.HTML(http.StatusOK, "main_nologged.html", c.Keys)
}

// ShowListenerTickets - показ билетов слушателя
func ShowListenerTickets(c *gin.Context) {
	loggedUser := GetLoggedUserFromSession(c)
	if !loggedUser.IsLogged {
		return
	}
	c.Set("LoggedUser", loggedUser)

	userTickets := mydatabase.FindUserLectionTicketsByField(loggedUser.User.ID, "", nil)
	c.Set("UserTickets", userTickets)
	c.HTML(http.StatusOK, "show_listener_tickets.html", c.Keys)
}

// LectureTickets - отображение занятых мест на лекцию
func LectureTickets(c *gin.Context) {
	lectureIDStr := c.Query("lecture_id")
	lectureID, err := strconv.Atoi(lectureIDStr)

	if err != nil {
		// если lecture_id - не число
		mysession.AddWarningFlash(c, "lecture_id")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	usersLecturesTickets := mydatabase.FindUserLectionTicketsByField(0, "lecture_id", lectureID)
	// panic(fmt.Sprintf("%v", usersLecturesTickets))
	c.Set("tickets", usersLecturesTickets)
	c.HTML(http.StatusOK, "lecture_tickets.html", c.Keys)
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

// ShowAllLectures - shows all lectures
func ShowAllLectures(c *gin.Context) {
	leclocs := mydatabase.FindLecturesLocationsByField("", "")
	c.Set("leclocs", leclocs)
	c.HTML(http.StatusOK, "show_all_lectures.html", c.Keys)
}

// BuyTicket - покупка пользователем билета на лекцию
func BuyTicket(c *gin.Context) {
	u := GetLoggedUserFromSession(c)
	lectureIDstr := c.Query("lecture_id")
	lectureID, err := strconv.Atoi(lectureIDstr)
	ifErr.Panic("lecture_id:", err)
	//  нет ли уже забронированных пользователем мест:
	userLectures := mydatabase.FindUserLectionTicketsByField(u.User.ID, "lecture_id", lectureID)
	if len(userLectures) > 0 {
		mysession.AddWarningFlash(c, "Вами уже забронировано место на это мероприятие")
		// редиректим назад
		c.Redirect(http.StatusTemporaryRedirect, c.Request.Referer())
		return
	}

	ok := mydatabase.BuyTicket(u.User.ID, lectureID)
	if !ok {
		SetWarningMsg(c, "ошибка, что-то пошло не так")
	}
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func ReleaseListenerTicket(c *gin.Context) {
	u := GetLoggedUserFromSession(c)
	ticketIDStr := c.PostForm("ticket_id")
	ticketID, err := strconv.Atoi(ticketIDStr)
	ifErr.Panic("can't convert ticket_id to int", err)
	ok := mydatabase.ReleaseTicket(ticketID, u.User.ID)
	if !ok {
		SetWarningMsg(c, "Что-то пошло не так")
	}
	c.Redirect(http.StatusTemporaryRedirect, "/")
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
func GetLoggedUserFromSession(c *gin.Context) LoggedUser {
	userFromSess, ok := GetUserFromSession(c)
	var lu = LoggedUser{
		IsLogged: ok,
		User:     userFromSess,
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

// LecturesOnLocation - shows lections on loation
func LecturesOnLocation(c *gin.Context) {
	locationIDStr := c.Query("location_id")
	c.Set("LocationID", locationIDStr)
	c.Set("lectures", mydatabase.FindLecturesByField("location_id", locationIDStr))
	c.HTML(http.StatusOK, "manage_lectures.html", c.Keys)
}

// InsertLecture - inserts lecture to location
func InsertLecture(c *gin.Context) {
	var lectureForm mydatabase.Lecture
	if err := c.ShouldBind(&lectureForm); err != nil {
		SetWarningMsg(c, "location_id error")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	// td: защититься от вставки к чужой площадке
	mydatabase.AddLecture(lectureForm)
	c.Redirect(http.StatusTemporaryRedirect, "/manage/lectures/?location_id="+strconv.Itoa(lectureForm.LocationID))
}

// SaveLecture - save lecture to location
func SaveLecture(c *gin.Context) {
	var lectureForm mydatabase.Lecture
	if err := c.ShouldBind(&lectureForm); err != nil {
		SetWarningMsg(c, "location_id error")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	// td: защититься от запроса к чужой площадке
	mydatabase.SaveLecture(lectureForm)
	c.Redirect(http.StatusTemporaryRedirect, "/manage/lectures/?location_id="+strconv.Itoa(lectureForm.LocationID))
}

// DeleteLecture - deletes lecture
func DeleteLecture(c *gin.Context) {
	locationID := c.Query("location_id")
	lectureID := c.Query("lecture_id")
	if len(lectureID) == 0 {
		SetWarningMsg(c, "lecture_id!")
		c.Redirect(http.StatusTemporaryRedirect, "/manage/lectures/?location_id="+locationID)
	}
	// если на лекцию кто-то зарегистрирован, то что делать? рассылать уведомления?
	// если билеты платные нужно деньги возвращать
	// пока просто запрещаем непустую лекцию удалять

	// найти слушателей:
	tickets := mydatabase.FindUserLectionTicketsByField(0, "lecture_id", lectureID)
	if len(tickets) > 0 {
		mysession.AddWarningFlash(c, "Чтоы удалить эту лекцию, удалите все билеты слушателей")
		c.Redirect(http.StatusTemporaryRedirect, "/manage/lectures/tickets/?lecture_id="+lectureID)
		return
	}

	mydatabase.DeleteLecture(lectureID)
	c.Redirect(http.StatusTemporaryRedirect, "/manage/lectures/?location_id="+locationID)
}

// EditLecture - show edit lecture page
func EditLecture(c *gin.Context) {
	var lectureID = c.Query("lecture_id")
	if len(lectureID) == 0 {
		SetWarningMsg(c, "lecture_id!")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	lectures := mydatabase.FindLecturesByField("id", lectureID)
	if len(lectures) == 0 {
		SetWarningMsg(c, "Not found lecture "+lectureID)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	c.Set("Lecture", lectures[0])
	c.HTML(http.StatusOK, "edit_lecture.html", c.Keys)
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

	// если не указана площадка спросить id площадки:
	if len(locID) == 0 {
		c.Set("locations", mydatabase.FindLocationsByField("", ""))
		c.HTML(http.StatusOK, "select_location.html", c.Keys)
		return
	}
	// если не указан юзер, спросить id юзера:
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

	// добавить юзеру роль организатора
	mydatabase.AddUserRole(userIDint, mydatabase.UserRoleOrganizer)
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
	// if len(userID) > 0 {
	locations := mydatabase.FindLocationsByField(userID, "")
	// }else{

	// }
	c.Set("locations", locations)
	// panic(fmt.Sprintf("%v", locations))
	// c.Set("locations", mydatabase.FindLocationsByField("", ""))
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
		mysession.AddWarningFlash(c, "Чтобы удалить площадку, удалите всех организаторов с данной лощадки")
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
	mydatabase.UpdateLocation(formLocation)
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
	if u.IsLogged && u.User.Roles >= 4 { // admin
		return
	}
	// если не админ:
	mysession.AddWarningFlash(c, "Недостаточно прав")
	c.Redirect(http.StatusTemporaryRedirect, "/")
	c.Abort()
}

// LoggedUser - IsLogged field answers on question "Is user logged?"
type LoggedUser struct {
	IsLogged bool
	User     mydatabase.User
}

// RegistrationEnd - обработка формы регистрации
func RegistrationEnd(c *gin.Context) {
	user := mydatabase.User{
		Email:    c.PostForm("email"),
		Password: GenerateSecret(),
		Fio:      c.PostForm("fio"),
		Roles:    1,
		ResetKey: "",
	}

	existingUser, ok := mydatabase.FindUserByEmail(user.Email)
	if ok {
		// Пользователь уже существует ( todo: добавить сброс пароля)
		// генерим ключ сброса
		existingUser.ResetKey = GenerateSecret() + GenerateSecret() // два подряд, и три влезло бы
		// сохраняем ключ в бд
		mydatabase.UpdateUser(existingUser)
		// отправляем ссылку на сброс
		// todo:
		// показываем сообщение
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
		ifErr.Panic("Can't add user ", err)

		mysession.AddInfoFlash(c, "Вам на почту отправлено письмо с паролем."+
			" (Если письмо долго не приходит - проверьте папку для спама).")
		c.Redirect(http.StatusTemporaryRedirect, "/login")

	}
}

// LoginEnd - обработка данных формы логина
func LoginEnd(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	user, ok := mydatabase.FindUserByEmail(email)
	if !ok {
		c.Set("warning_msg", "Пользователь с таким адресом не найден. "+
			"Попробуйте получить пароль.")
		c.HTML(http.StatusOK, "login.html", c.Keys)
	} else {
		if password != user.Password {
			// юзер найден но пароль не совпадает:
			c.Set("warning_msg", "Пароль почему-то не тот, который мы Вам высылали.")
			c.HTML(http.StatusOK, "login.html", c.Keys)
		} else {
			// юзер найден и емаил подтвержден:
			mysession.SaveEmail(c, email)
			// c.HTML(http.StatusOK, "main.html", gin.H{})
			// ShowMainPage(c)
			c.Redirect(http.StatusTemporaryRedirect, "/")
		}
	}
}
