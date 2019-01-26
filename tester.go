package main

import (
	"fmt"
	"net/http"
	"net/mail"
	"os"

	"github.com/olegbelyaev/siteforreg/app"
	"github.com/olegbelyaev/siteforreg/myemail"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/olegbelyaev/siteforreg/mydatabase"
)

func main() {

	// добавление суперадмина
	mydatabase.AddInitAdmin()

	// конфигурация для отправки почты:
	myemail.SetParams(
		"", "sivsite@yandex.ru", os.Getenv("EMAIL_SECRET"),
		"smtp.yandex.ru", "465",
		mail.Address{Name: "sitename", Address: "sivsite@yandex.ru"},
		true,
	)

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.Use(func(c *gin.Context) {
		c.Set("html_title", "Siteforrreg")
		c.Set("LoggedUser", app.GetLoggedUserFromSession(c))
	})

	// ======================== главная / регистрация / логин / выход =====================

	// ------------------ маршрут к главной странице -----------------
	router.GET("/", app.ShowMainPage)

	// ----------------- форма логина -----------------------------
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", c.Keys)
	})

	// ---------------- обработка данных формы логина -----------------
	router.POST("/login/end", app.LoginEnd)

	// --------------------- разлогинивание --------------------------
	router.GET("/logout", func(c *gin.Context) {
		// для разлогина сохраним емаил, по которому пользователь не найдется
		app.SaveEmailToSession(c, "")
		app.ShowMainPage(c)
	})

	// ----------------- форма регистрации на сайте ----------------------
	router.GET("/registration/start", func(c *gin.Context) {
		c.HTML(http.StatusOK, "registration.html", c.Keys)
	})

	// ----------- обработка формы регистрации нового пользователя ------------
	router.POST("/registration/end", app.RegistrationEnd)

	// ================== площадки =========================================

	locations := router.Group("/locations")
	{
		// ------------ список площадок -----------------------------------
		locations.Any("/", app.ShowLocations)

		// ниже этого будет требовать залогиниться:
		locations.Use(app.GotoLoginIfNotLogged)
		// требование быть админом:
		locations.Use(app.GotoAccessDeniedIfNotAdmin)

		// ------------- форма добавления новой площадки ---------------------
		locations.GET("/new", func(c *gin.Context) {
			c.HTML(http.StatusOK, "tmp_valid_locations.html", c.Keys)
		})

		// ------------- обработка формы добавления новой площадки ----------
		locations.POST("/insert", app.AddLocation)

		// td: как защититься от запросов не с этого сайта?
		// ----------- обработка запроса от формы удаления площадки ------------
		locations.Any("/delete/:ID", app.DeleteLocation)
	}

	// ========================== пользователи ================================
	users := router.Group("/users")
	{
		users.Use(app.GotoLoginIfNotLogged)
		users.Use(app.GotoAccessDeniedIfNotAdmin)
		// -------------- список ------------------
		users.GET("/", app.ShowUsers)
	}

	// ======================= организаторы на площадках ===================
	// ------------- список --------------------------------------------
	router.GET("/locorgs", app.ShowLocorgs)

	// ------------- добавление организатора на площадку ----------------
	router.Any("/add_locorg", app.AddLocOrg)

	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> запуск! <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	port := os.Getenv("PORT")
	if len(port) == 0 {
		panic("PORT env not defined!")
	}

	router.Run(fmt.Sprintf(":%s", port))
}
