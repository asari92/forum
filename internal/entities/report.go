package entities

type Report struct {
	ID           int
	UserID       int
	ReporterName string
	PostID       int
	Reason       string
	Created      string
}
