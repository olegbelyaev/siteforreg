package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olegbelyaev/siteforreg/mydatabase"
)

// ShowUsers -список пользователей
func ShowUsers(c *gin.Context) {
	c.Set("users", mydatabase.FindUsersByField("", ""))
	c.HTML(http.StatusOK, "show_users.html", c.Keys)
}
