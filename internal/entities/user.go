package entities

type User struct {
	ID             int
	Username       string
	Email          string
	HashedPassword []byte
	Created        string
	Role           string // "user", "moderator", "admin"
}

const (
	RoleUser      = "user"
	RoleModerator = "moderator"
	RoleAdmin     = "admin"
)
