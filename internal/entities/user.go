package entities

type User struct {
	ID             int
	Username       string
	Email          string
	HashedPassword []byte
	Created        string
}