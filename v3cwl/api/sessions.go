package api


import (
	"fmt"
	"log"
	"path"
	"net/http"

	"cwl/types"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/gorilla/sessions"
)




type SessionOptions struct {

	CookiesKey 	string
	MaxAge     	int
	HttpOnly 	bool
	Secure   	bool
}



func NewCookieStore(opts SessionOptions) *sessions.CookieStore {

	store := sessions.NewCookieStore([]byte(opts.CookiesKey))

	store.MaxAge(opts.MaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = opts.HttpOnly
	store.Options.Secure = opts.Secure

	return store
}




func (s *AuthService) StoreUserSession(w http.ResponseWriter, r *http.Request, user goth.User) error {

	avatarHash := path.Base(user.AvatarURL)
	nickname := user.RawData["global_name"].(string)

	id, err := s.store.Users.Upsert(&types.DiscordUser{

		DiscordID: user.UserID,
		Nickname: nickname,
		Username: user.Name,
		AvatarHash: avatarHash,
	})
	
	if err != nil {
		return err
	}

	session, _ := gothic.Store.Get(r, "user-session")

	session.Values["user_id"] = id

	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}



type sessionData struct {
	UserID int
}

func (s *AuthService) GetUserSession(r *http.Request) (sessionData, error) {

	session, err := gothic.Store.Get(r, "user-session")
	if err != nil {
		return sessionData{}, err
	}

	id := session.Values["user_id"]

	if id == nil {
		return sessionData{}, fmt.Errorf("user is not authenticated: %v", id)
	}

	
	sessionData := sessionData{
		UserID: id.(int),
	}

	return sessionData, nil
}




func (s *AuthService) ClearUserSession(w http.ResponseWriter, r *http.Request) {

	session, err := gothic.Store.Get(r, "user-session")

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Clearing session: %v\n", session.Values["user_id"])
	
	session.Values["user_id"] = nil
	session.Options.MaxAge = -1

	err = session.Save(r, w)
}