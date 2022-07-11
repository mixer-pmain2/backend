package session

import (
	"github.com/gorilla/sessions"
	"net/http"
)

var (
	key   = []byte("secret-key")
	Store = sessions.NewCookieStore(key)
)

type sessionStore struct {
	w    http.ResponseWriter
	r    *http.Request
	name string
}

type user struct {
	sessionStore
	Lname  string
	Fname  string
	Sname  string
	IsAuth bool
	Id     int64
}

func GetUserStore(w http.ResponseWriter, r *http.Request) *user {
	session, _ := Store.Get(r, "user")
	data := user{
		sessionStore: sessionStore{w, r, "user"},
		IsAuth:       false,
		Id:           0,
	}
	data.Lname, _ = session.Values["lname"].(string)
	data.Fname, _ = session.Values["fname"].(string)
	data.Sname, _ = session.Values["sname"].(string)
	data.IsAuth, _ = session.Values["isAuth"].(bool)
	data.Id, _ = session.Values["id"].(int64)

	return &data
}

func (u *user) Save() {
	session, _ := Store.Get(u.r, u.name)
	session.Values["lname"] = u.Lname
	session.Values["fname"] = u.Fname
	session.Values["sname"] = u.Sname
	session.Values["isAuth"] = u.IsAuth
	session.Values["id"] = u.Id
	session.Save(u.r, u.w)
}
