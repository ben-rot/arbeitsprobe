package types 


type FamilyProfile struct {
	
	DiscordServerID 	string		`json:"-"   db:"discord_server_id"`
	Name 				string 		`json:"-"   db:"name"`
	IconHash 			string 		`json:"-"   db:"server_icon"`
}



type Family struct {

	ID 					int			`json:"-"   db:"id"`
	FamilyProfile

	Managers 			[]int		`json:"-"   db:"-"` // User IDs
	Clans 				[]Clan		`json:"-"   db:"-"` // Clan
	Accounts 			[]Account	`json:"-"   db:"-"` // Account
}

// 3 Managers Max


type FamilyStats struct {
	FamilyProfile
	TotalAccounts 		int	 		`json:"-" db:"total_accounts"`
	AssignedAccounts 	int
}




func (f *Family) SortPlayers() {

	clanMap := make(map[int][]Account)
	var unassigned []Account


	for _, account := range f.Accounts {

		if account.ClanID == 0 {
			unassigned = append(unassigned, account)
		} else {
			clanMap[account.ClanID] = append(clanMap[account.ClanID], account)
		}
	}

	for i := range f.Clans {

		clanId := f.Clans[i].ID

		if accs, ok := clanMap[clanId]; ok {
			f.Clans[i].Accounts = accs
			delete(clanMap, clanId)
		}
	}

	f.Accounts = unassigned
}