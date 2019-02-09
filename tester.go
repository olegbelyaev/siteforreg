package main

import (
	"fmt"
	"net/http"
	"net/mail"
	"os"

	"github.com/olegbelyaev/siteforreg/app"
	"github.com/olegbelyaev/siteforreg/myemail"
	"github.com/olegbelyaev/siteforreg/mysession"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/olegbelyaev/siteforreg/mydatabase"
)

func main() {

	// добавление суперадмина
	mydatabase.AddInitAdmin()
	mydatabase.AddInitRoles()

	// конфигурация для отправки почты:
	myemail.SetParams(
		"", "sivsite@yandex.ru", os.Getenv("EMAIL_SECRET"),
		"smtp.yandex.ru", "465",
		mail.Address{Name: "MeetFor", Address: "sivsite@yandex.ru"},
		true,
	)

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.Use(func(c *gin.Context) {
		c.Set("html_title", "Meet for")
		c.Set("LoggedUser", app.GetLoggedUserFromSession(c))
		warningFlashes, infoFlashes := mysession.GetFlashes(c)
		c.Set("WarningFlashes", warningFlashes)
		c.Set("InfoFlashes", infoFlashes)
	})

	// ======================== главная / регистрация / логин / выход =====================

	router.Any("/", app.ShowMainPage)

	router.Any("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", c.Keys)
	})

	router.POST("/login/end", app.LoginEnd)

	router.GET("/logout", func(c *gin.Context) {
		// для разлогина сохраним емаил, по которому пользователь не найдется
		mysession.SaveEmail(c, "")
		app.ShowMainPage(c)
	})

	router.GET("/registration/start", func(c *gin.Context) {
		c.HTML(http.StatusOK, "registration.html", c.Keys)
	})

	router.POST("/registration/end", app.RegistrationEnd)

	// td: как защититься от запросов не с этого сайта?
	// зона администратора:
	administrate := router.Group("/administrate")
	{
		administrate.Use(app.GotoLoginIfNotLogged)
		administrate.Use(app.GotoAccessDeniedIfNotAdmin)

		// ================== площадки =========================================
		locations := administrate.Group("/locations")
		{
			locations.Any("/", app.ShowLocations)

			locations.GET("/new", func(c *gin.Context) {
				c.HTML(http.StatusOK, "new_location_form.html", c.Keys)
			})

			locations.POST("/insert", app.InsertLocation)

			locations.Any("/delete/:ID", app.DeleteLocation)

			locations.Any("/edit/:ID", app.EditLocation)
			locations.POST("/save", app.SaveLocation)
		}

		// ========================== пользователи ================================
		users := administrate.Group("/users")
		{
			users.GET("/", app.ShowUsers)
		}

		// ======================= организаторы на площадках ===================
		locorgs := administrate.Group("/locorgs")
		{
			locorgs.GET("/", app.ShowLocorgs)

			locorgs.Use(app.GotoLoginIfNotLogged)
			locorgs.Use(app.GotoAccessDeniedIfNotAdmin)

			locorgs.Any("/add_locorg", app.AddLocOrg)
			locorgs.Any("/delete", app.DeleteLocorg)
		}
	}

	manage := router.Group("/manage")
	{
		manage.Use(app.GotoLoginIfNotLogged)

		// ======================== площадки юзера-организатора ====================
		mylocorgs := manage.Group("/mylocorgs")
		{
			mylocorgs.Any("/", app.ShowMyLocOrgs)

		}

		mylectures := manage.Group("/lectures")
		{
			mylectures.Any("/", app.ShowMyLectures)

			mylectures.Any("/new", func(c *gin.Context) {
				LocationID := c.Query("location_id")
				if len("location_id") == 0 {
					app.SetWarningMsg(c, "location_id not defined")
					c.Redirect(http.StatusTemporaryRedirect, "/manage/lectures/")
					return
				}
				c.Set("LocationID", LocationID)
				c.HTML(http.StatusOK, "new_lecture.html", c.Keys)
			})
			mylectures.POST("/insert", app.InsertLecture)

			mylectures.Any("edit", app.EditLecture)
			mylectures.POST("/save", app.SaveLecture)
			mylectures.Any("/delete", app.DeleteLecture)
		}
	}

	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> запуск! <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	port := os.Getenv("PORT")
	if len(port) == 0 {
		panic("PORT env not defined!")
	}

	router.Run(fmt.Sprintf(":%s", port))
}
