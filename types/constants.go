package types


type ContextKey string
const UserIdKey ContextKey = "userId"



type URLFrag string
const (
	DiscordAPI URLFrag = "https://discord.com/api/v9/"
	DiscordCDN URLFrag = "https://cdn.discordapp.com/"
)