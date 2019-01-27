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

// SetWarningMsg - установка сообщения, которое пробросится в шаблоны ключем "warning_msg"
func SetWarningMsg(c *gin.Context, msg string) {
	c.Set("warning_msg", msg)
}
