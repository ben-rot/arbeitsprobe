package types 


type DiscordUser struct {

	DiscordID 		string 		`db:"discord_id"`
	Username 		string 		`db:"username"`
	Nickname 		string 		`db:"nickname"`
	AvatarHash 		string 		`db:"avatar_hash"`
}


type User struct {

	ID 				int 			`db:"id"`
	DiscordUser
}