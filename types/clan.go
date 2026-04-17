package types 


type ClanProfile struct {
	Tag 			string     	`json:"tag"             db:"tag"`
	Name 			string		`json:"name"            db:"name"`
	League 			int			`json:"warLeague"       db:"cwl_league"`
	IsLocked 		bool		`json:"-"               db:"is_locked"`
}

type Clan struct {

	ID 				int			`json:"-"       db:"id"`
	ClanProfile

	Accounts 		[]Account
}



type ClanStats struct {
	ClanProfile
	AccountCount 	int 	`json:"-" db:"acc_count"`
}