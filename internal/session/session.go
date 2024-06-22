package session

import (
	"fmt"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/gorilla/sessions"
	"net/http"
	"os"
)

var (
	key               = []byte(common.Ternary(os.Getenv("SESSION_KEY") != "", os.Getenv("SESSION_KEY"), "secret").(string))
	store             = sessions.NewCookieStore(key)
	sessionCookieName = "user-session" // using this name to store session cookie
)

func CreateUserSession(writer http.ResponseWriter, request *http.Request) error {
	session, _ := store.Get(request, sessionCookieName)
	session.Values["authenticated"] = true
	return session.Save(request, writer)
}

func IsAuthenticated(request *http.Request) bool {
	session, err := store.Get(request, sessionCookieName)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if session.IsNew {
		return false
	}
	return session.Values["authenticated"].(bool)
}

func Logout(writer http.ResponseWriter, request *http.Request) error {
	session, _ := store.Get(request, sessionCookieName)
	session.Values["authenticated"] = false
	return session.Save(request, writer)
}
