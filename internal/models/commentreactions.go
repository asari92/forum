package models

import "database/sql"

type CommentReaction struct {
	UserID    int
	CommentID int
	IsLike    bool
}

type CommentReactionModel struct {
	DB *sql.DB
}
