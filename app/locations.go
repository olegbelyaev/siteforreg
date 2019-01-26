package app

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/olegbelyaev/siteforreg/mydatabase"
)

// ShowLocations - список площадок
func ShowLocations(c *gin.Context) {
	userID := c.Query("user_id")
	if len(userID) > 0 {
		c.Set("locations", mydatabase.FindLocationsByField(userID, ""))
	}
	c.Set("locations", mydatabase.FindLocationsByField("", ""))
	c.HTML(http.StatusOK, "locations.html", c.Keys)
}

// AddLocation - добавление новой площадки
func AddLocation(c *gin.Context) {
	var l mydatabase.Location
	// получение данных из формы создания новой площадки:
	if err := c.ShouldBind(&l); err != nil {
		c.Set("warning_msg", err.Error())
		c.HTML(http.StatusOK, "templateAddLocation.html", c.Keys)
	}
	mydatabase.AddLocation(l)
	c.Redirect(http.StatusTemporaryRedirect, "/locations/")
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

	err = mydatabase.DeleteLocation(locationID)
	if err != nil {
		c.Set("warning_msg", err.Error())
	}
	ShowLocations(c)
}
