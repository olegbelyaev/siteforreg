package mysession

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// SessionName - default session name for this app
const SessionName = "site-forreg"

var sessionSecret string

// CookieStore - returns singleton cookieStore,
// created by secret from "SESSION_SECRET" env var
func CookieStore() *sessions.CookieStore {
	if len(sessionSecret) == 0 {
		sessionSecret = strings.TrimSpace(os.Getenv("SESSION_SECRET"))
	}
	if len(sessionSecret) == 0 {
		panic("SESSION_SECRET env var not defined. Set up it and restart application.")
	}

	cookieStore := sessions.NewCookieStore([]byte(sessionSecret))
	// cookieStore := sessions.NewCookieStore([]byte("SESSION_SECRET"))
	return cookieStore
}

// GetSession - get cookie store session
func GetSession(c *gin.Context) *sessions.Session {
	store := CookieStore()
	sess, _ := store.Get(c.Request, SessionName)
	// if err != nil {
	// 	panic("Can't get session: " + err.Error() + " " + string(sessionSecret))
	// }
	return sess
}

// GetStringValue - возвращает ранее сохраненное в сессии строковое значение по ключу key
func GetStringValue(c *gin.Context, key string) (string, bool) {
	sess := GetSession(c)
	// panic(sess.Values)
	email, ok := sess.Values[key]
	// panic("email:" + email.(string))
	if !ok {
		return "", ok
	}
	return email.(string), ok
}
