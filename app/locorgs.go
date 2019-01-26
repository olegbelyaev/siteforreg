package app

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/olegbelyaev/siteforreg/mydatabase"
)

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

// AddLocOrg - организаторы на площадках
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

}
