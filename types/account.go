package types 




type ClashData struct {
	Tag 			string		`json:"-"               db:"tag"`
	Name 			string		`json:"name"            db:"name"`
	Townhall 		int			`json:"townHallLevel"   db:"townhall"`
}


type Account struct {

	ID 				int 		`json:"-"               db:"id"`
	ClashData

	OwnerID 		int			`json:"-"               db:"owner_id"`
	FamilyID 		int			`json:"-"               db:"family_id"`
	ClanID 			int			`json:"-"               db:"clan_id"`
}



type DashboardAccount struct {
	ID 				int 		`db:"id"`
	ClashData

	FamilyName 		*string 	`db:"family_name"`

	ClanTag 		*string 	`db:"clan_tag"`
	ClanName 		*string 	`db:"clan_name"`
	ClanLeague 		*int 		`db:"clan_league"`
}