package repositories

import "forum/entities"

type CommentReactionRepository interface {
	AddReaction(userID, commentID int, isLike bool) error
	RemoveReaction(userID, commentID int) error
	GetLikesCount(commentID int) (int, error)
	GetDislikesCount(commentID int) (int, error)
	GetUserReaction(userID, commentID int) (*entities.CommentReaction, error)
}
