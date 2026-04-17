package api


import (
	"fmt"
	"log"
	"net/http"

	"cwl/config"
	"cwl/database"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
)








type AuthService struct {
	sessionStore sessions.Store
	store 		*database.Store
}


func NewAuthService(store sessions.Store, authCfg *config.AuthConfig, db *database.Store) *AuthService {

	gothic.Store = store

	goth.UseProviders(
		
		discord.New(
			authCfg.ClientId,
			authCfg.ClientSecret,
			authCfg.CallbackUrl,
			"identify guilds",
		),
	)

	return &AuthService{
		sessionStore: 	store,
		store: 			db,
	}
}












func (s *AuthService) Login(w http.ResponseWriter, r *http.Request) {

	id, err := s.GetUserSession(r)

	if err == nil {

		log.Printf("User is already authenticated: %v", id)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	} 

	gothic.BeginAuthHandler(w, r)	
}





func (s *AuthService) LoginCallback(w http.ResponseWriter, r *http.Request) {

	u, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	err = s.StoreUserSession(w, r, u)
	if err != nil {
		log.Println(err)
		log.Println("Failed to store user session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	log.Printf("User is authenticated: %v", u.Name)

	http.Redirect(w, r, "/home", http.StatusTemporaryRedirect)
}






func (s *AuthService) Logout(w http.ResponseWriter, r *http.Request) {

	err := gothic.Logout(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	w.Header().Set("HX-Redirect", "/")


	s.ClearUserSession(w, r)

	w.WriteHeader(http.StatusOK)
}