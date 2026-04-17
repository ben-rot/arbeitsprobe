package types

import "github.com/a-h/templ"


type PageData struct {

	Title 		string
	Indicator 	templ.Component
	IsManager 	bool
}



type HomeData struct {
	User 			*User
	Accounts 		[]DashboardAccount
	FamilyStats 	*FamilyStats
	ClanStats 		[]ClanStats
}