package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/a-h/templ"

	"cwl/config"
	"cwl/database"
	"cwl/views"
	"cwl/types"
)



type Server struct {
	listenAddr 	string
	cfg        	*config.Config
	store       *database.Store
	auth		*AuthService
	clashClient *ClashClient
}




func NewServer(listenAddr string, cfg *config.Config, store *database.Store, clashClient *ClashClient) *Server {

	sessionStore := NewCookieStore(SessionOptions{
		CookiesKey: 	cfg.Auth.ClientSecret,
		MaxAge:     	86400,
		HttpOnly: 		true, 
		Secure:   		false,
	})

	authService := NewAuthService(sessionStore, &cfg.Auth, store)

	return &Server{
		listenAddr: 	listenAddr,
		cfg:        	cfg,
		store:         	store,
		auth: 			authService,
		clashClient: 	clashClient,
	}
}




func (s *Server) render(w http.ResponseWriter, r *http.Request, pageData *types.PageData, content templ.Component) {

	isHTMX := r.Header.Get("HX-Request") == "true"

	if isHTMX {
		content.Render(r.Context(), w)
		return
	}


	if !pageData.IsManager {

		session, err := s.auth.GetUserSession(r)
		if err != nil { return }

		managerFamilyId, err := s.store.Users.GetManagedFamilyId(session.UserID)
		if err != nil { return }

		pageData.IsManager = managerFamilyId != nil
	}
	

	views.Page(*pageData).Render(templ.WithChildren(r.Context(), content), w)
}




func (s *Server) Start() error {

	router := mux.NewRouter().StrictSlash(true)

	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	router.HandleFunc("/", s.handleIndex)
	s.registerRoutes(router)

	return http.ListenAndServe(s.listenAddr, router)
}