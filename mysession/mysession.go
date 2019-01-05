package mysession

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// SessionName - default session name for this app
const SessionName = "site-forreg"

var cookieStore *sessions.CookieStore

// CookieStore - returns singleton cookieStore,
// created by secret from "SESSION_SECRET" env var
func CookieStore() *sessions.CookieStore {
	if cookieStore == nil {
		secr := os.Getenv("SESSION_SECRET")
		if len(secr) == 0 {
			panic("SESSION_SECRET env var not defined. Set up it and restart application.")
		}
		// TODO: не работает env
		// cookieStore = sessions.NewCookieStore([]byte(secr))
		cookieStore = sessions.NewCookieStore([]byte("SESSION_SECRET"))
	}
	return cookieStore
}

// GetSession - get cookie store session
func GetSession(c *gin.Context) *sessions.Session {
	store := CookieStore()
	sess, err := store.Get(c.Request, SessionName)
	if err != nil {
		panic("Can't get session: " + err.Error())
	}
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
