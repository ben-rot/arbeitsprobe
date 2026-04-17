package api

import (
	"context"
	"log"
	"fmt"
	"net/http"

	"cwl/types"

	"github.com/gorilla/mux"
)


func (s *Server) RequireAuth(next http.Handler) http.Handler {

	return http.HandlerFunc( 

		func(w http.ResponseWriter, r *http.Request) {

			sessionData, err := s.auth.GetUserSession(r)

			if err != nil {
				log.Println("User is not authenticated, redirecting to login")

				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Redirect", "/")
					w.WriteHeader(http.StatusOK)
					return
				}

				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return 
			}

			ctx := context.WithValue(r.Context(), types.UserIdKey, sessionData.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}


// duplicatePageClickPrevention 



// In progress
func (s *Server) RequireManager(next http.Handler) http.Handler {


	return http.HandlerFunc(

		func(w http.ResponseWriter, r *http.Request) {

			userId, ok := r.Context().Value("userId").(int)
			if !ok {
				fmt.Println("Could not get userId")
				return
			}

			familyId, err := s.store.Users.GetManagedFamilyId(userId)
			if err != nil {
				fmt.Println("Could not get managed family id: ", err)
				return
			}

			if familyId == nil {
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}

			ctx := context.WithValue(r.Context(), "familyId", *familyId)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}





func (s *Server) ConvertFamilyIdParam(next http.Handler) http.Handler {

	return http.HandlerFunc(
		
		func (w http.ResponseWriter, r *http.Request) {

			discordId := mux.Vars(r)["familyId"]
			fmt.Println("Discord ID: ", discordId)

			query := `SELECT id FROM family WHERE discord_server_id = $1`
			var familyId int
			err := s.store.DB.Get(&familyId, query, discordId)
			if err != nil {
				fmt.Println("Error:", err)
				http.Error(w, "Error:" + err.Error(), http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), "familyId", familyId)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}