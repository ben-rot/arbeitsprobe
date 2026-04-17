package api

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	"cwl/views"
	"cwl/types"
)





func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Cache-Control", "no-store, must-revalidate")

	_, err := s.auth.GetUserSession(r)

	if err != nil {

		views.Login().Render(r.Context(), w)
		return
	}
	
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Location", "/home")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}






func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Routed: Home")

	userId, ok := r.Context().Value(types.UserIdKey).(int)
	if !ok {
		fmt.Println("Could not get userId")
		return
	}


	data := types.HomeData{}

	var (
		wg 					sync.WaitGroup
		managedFamilyId 	*int
		errs = make(chan error, 3)
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		var err error
		data.User, err = s.store.Users.GetById(userId)
		errs <- err
	}()

	go func() {
		defer wg.Done()
		var err error
		data.Accounts, err = s.store.Accounts.DashboardDataByOwnerId(userId)
		errs <- err
	}()

	go func() {
		defer wg.Done()
		var err error
		managedFamilyId, err = s.store.Users.GetManagedFamilyId(userId)
		errs <- err
	}()

	wg.Wait()
	close(errs)


	for err := range errs {
        if err != nil {
            fmt.Println("Parallel fetch error: ", err)
            return
        }
    }


	if managedFamilyId != nil {

		var (
			wg 			sync.WaitGroup
			errs = make(chan error, 2)
		)
		wg.Add(2)

		go func() {
			defer wg.Done()
			var err error
			data.FamilyStats, err = s.store.Families.GetStatsForDashboard(*managedFamilyId)
			errs <- err
		}()

		go func() {
			defer wg.Done()
			var err error
			data.ClanStats, err = s.store.Clans.GetStatsForDashboard(*managedFamilyId)
			errs <- err
		}()

		wg.Wait()
		close(errs)

		for err := range errs {
			if err != nil {
				fmt.Println("Could not get stats: ", err)
				return
			}
		}
	}


		
	pageData := &types.PageData{
		Title: "Home",
		Indicator: views.GenericIndicator("Home"),
		IsManager: managedFamilyId != nil,
	}

	s.render(w, r, pageData, views.Home(data))
}







func (s *Server) handleAccounts(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Routed: Accounts")

	userId, ok := r.Context().Value(types.UserIdKey).(int)
	if !ok {
		fmt.Println("Could not get userId")
		return
	}



	accounts, err := s.store.Accounts.OfOwnerById(userId)
	if err != nil {
		fmt.Println("Could not get accounts: ", err)
		return
	}



	pageData := &types.PageData{
		Title: "My Accounts",
		Indicator: views.GenericIndicator("Accounts"),
	}

	s.render(w, r, pageData, views.Accounts(accounts))
}






func (s *Server) handleBoard(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Routed: Studio")


	pageData := &types.PageData{
		Title: "Studio",
		Indicator: views.StudioIndicator("", "", ""),
	}


	s.render(w, r, pageData, views.Editor(0))
}





















func (s *Server) handleAddAccountForm(w http.ResponseWriter, r *http.Request) {

	views.AddAccountModal().Render(r.Context(), w)
}




func (s *Server) registerNewAccount(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Registering new account")

	tag := s.store.Accounts.SanitizeTag(r.FormValue("account_tag"))
	fmt.Println(tag)


	// Check Duplicate in DB
	// Implement Life check later on 

	exists, err := s.store.Accounts.Exists(tag)

	if err != nil {
		fmt.Println(err)
		return
	}

	if exists {
		fmt.Println("Tag already exists")
		return
	}

	fmt.Println("Tag is unqiue")




	// Ownership Validation

	token := r.FormValue("api_token")
	
	isValid, err := s.clashClient.verifyAccount(tag, token)

	if err != nil || !isValid {
		fmt.Println(err)
		return
	}




	// Get Account Details

	id, ok := r.Context().Value(types.UserIdKey).(int)
	fmt.Println(id)
	if !ok {
		fmt.Println("Could not get userId")
		return
	}

	account, err := s.clashClient.getAccount(tag)

	if err != nil {
		fmt.Println("Could not get account")
		return
	}

	newAccount := types.Account{
		OwnerID: id,
		ClashData: *account,
	}

	err = s.store.Accounts.Save(&newAccount)
	if err != nil {
		fmt.Println("Could not save account: ", err)
		return
	}

	w.Header().Set("HX-Trigger", "clear-form")
	views.AccountCard(newAccount).Render(r.Context(), w)
	views.OobDelPlaceholder("acc-list-placeholder").Render(r.Context(), w)
	views.AddAccBtnOob(5).Render(r.Context(), w)
}





func (s *Server) deleteAccount(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Deleting account")
	tag := mux.Vars(r)["tag"]

	fmt.Println(tag)
	err := s.store.Accounts.Delete(tag)

	if err != nil {
		fmt.Println("Could not delete account")
		return
	}

	// return confirmation
}