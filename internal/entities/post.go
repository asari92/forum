package entities

type Post struct {
	ID         int
	Title      string
	Content    string
	UserID     int
	UserName   string
	Created    string
	IsApproved bool
}
