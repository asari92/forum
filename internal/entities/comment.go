package entities

type Comment struct {
	ID           int
	PostID       int
	UserID       int
	UserName     string
	UserRole string
	Content      string
	UserReaction int
	Like         int
	Dislike      int
	Created      string
}
