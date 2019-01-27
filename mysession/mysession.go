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
	return cookieStore
}

// GetSession - get cookie store session
// Note: если получен err!=nil, скорее всего это значит изменился ключ шифрования.
// Сессия создастся и будет возвращена новая.
// В этом случае можно игнорировать ошибку.
func GetSession(c *gin.Context) (*sessions.Session, error) {
	store := CookieStore()
	sess, err := store.Get(c.Request, SessionName)
	return sess, err
}

// SaveEmail  - сохраняет email пользователя в сессию
func SaveEmail(c *gin.Context, email string) {
	sess, _ := GetSession(c)
	sess.Values["email"] = email
	sess.Save(c.Request, c.Writer)
}

// GetStringValue - возвращает ранее сохраненное в сессии строковое значение по ключу key
func GetStringValue(c *gin.Context, key string) (string, bool) {
	sess, _ := GetSession(c)
	email, ok := sess.Values[key]
	if !ok {
		return "", ok
	}
	return email.(string), ok
}

// AddWarningFlash - добавить warning-флеш-сообщение пользователю
func AddWarningFlash(c *gin.Context, msg string) {
	sess, _ := GetSession(c)
	sess.AddFlash("WARNING:" + msg)
	sess.Save(c.Request, c.Writer)
}

// AddInfoFlash - добавить info-флеш-сообщение пользователю
func AddInfoFlash(c *gin.Context, msg string) {
	sess, _ := GetSession(c)
	sess.AddFlash("INFO:" + msg)
	sess.Save(c.Request, c.Writer)
}

// GetFlashes - возвращает список предупреждающих и иформационных флеш-сообщений (и удаляет их из сессии)
func GetFlashes(c *gin.Context) ([]string, []string) {
	sess, _ := GetSession(c)
	allFlashes := sess.Flashes()
	sess.Save(c.Request, c.Writer)
	var warningFlashes []string
	var infoFlashes []string
	for _, fl := range allFlashes {
		if strings.HasPrefix(fl.(string), "WARNING") {
			warningFlashes = append(warningFlashes, fl.(string))
			continue
		}
		infoFlashes = append(infoFlashes, fl.(string))
	}

	return warningFlashes, infoFlashes
}
