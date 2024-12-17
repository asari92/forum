package entities

type Notification struct {
	ID              int
	OwnerID         int
	PostID          int
	PostTitle       string
	PostContent     string
	Action          string
	Created         string
	TriggerUserID   int
	TriggerUserName string
}
