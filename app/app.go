package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ShowMainPage - показ главной страницы сайта
func ShowMainPage(c *gin.Context) {
	c.Set("LoggedUser", GetLoggedUserFromSession(c))
	c.HTML(http.StatusOK, "main.html", c.Keys)
}
