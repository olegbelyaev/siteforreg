package app

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/olegbelyaev/siteforreg/mydatabase"
	"github.com/olegbelyaev/siteforreg/mysession"
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
	// c.Redirect(http.StatusTemporaryRedirect, "/locations/")
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
		c.Redirect(http.StatusTemporaryRedirect, "/locorgs/?location_id="+locationIDStr)
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
