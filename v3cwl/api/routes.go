package api

import (
	"github.com/gorilla/mux"
)


func (s *Server) registerRoutes(r *mux.Router) {

	// Auth routes
	r.HandleFunc("/login/{provider}", s.auth.Login).Methods("GET")
	r.HandleFunc("/login/{provider}/callback", s.auth.LoginCallback).Methods("GET")
	r.HandleFunc("/logout", s.auth.Logout).Methods("GET")



	// view routes
	views := r.PathPrefix("/").Subrouter()
	views.HandleFunc("", s.handleIndex).Methods("GET")

	protectedViews := views.PathPrefix("/").Subrouter()
	protectedViews.Use(s.RequireAuth)

	protectedViews.HandleFunc("/studio", s.handleBoard).Methods("GET")
	protectedViews.HandleFunc("/home", s.handleHome).Methods("GET")
	protectedViews.HandleFunc("/accounts", s.handleAccounts).Methods("GET")



	// Accounts routes
	accountGroup := r.PathPrefix("/accounts").Subrouter()
	accountGroup.Use(s.RequireAuth)

	accountGroup.HandleFunc("", s.handleAccounts).Methods("GET")
	accountGroup.HandleFunc("/add", s.registerNewAccount).Methods("POST")
	accountGroup.HandleFunc("/add-form", s.handleAddAccountForm).Methods("GET")
	accountGroup.HandleFunc("/{tag}", s.deleteAccount).Methods("DELETE")


	// Family routes
	familyGroup := r.PathPrefix("/studio").Subrouter()
	familyGroup.Use(s.RequireAuth)
	// familyGroup.Use(s.ConvertFamilyIdParam)
	// familyGroup.Use(RequireManager)

	familyGroup.HandleFunc("", s.handleBoard).Methods("GET")
}