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

func (c *CommentReactionModel) AddReaction(userId, commentID int, isLike bool) error {
	stmt := `INSERT INTO comment_reactions (user_id, comment_id,is_like)
	VALUES (?,?,?)
	ON CONFLICT(user_id, comment_id) DO UPDATE SET is_like = ?`
	_, err := c.DB.Exec(stmt, userId, commentID, isLike, isLike)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommentReactionModel) RemoveReaction(userId, commentId int) error {
	stmt := `DELETE FROM comment_reactions
	WHERE user_id = ? AND comment_id = ?`
	_, err := c.DB.Exec(stmt, userId, commentId)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommentReactionModel) GetLikesCount(commentId int) (int, error) {
	stmt := `SELECT COUNT(*) FROM comment_reactions
	WHERE comment_id = ? AND is_like = 1`
	var likes int
	row := c.DB.QueryRow(stmt, commentId)
	err := row.Scan(&likes)
	if err != nil {
		return 0, err
	}
	return likes, nil
}

func (c *CommentReactionModel) GetDislikesCount(commentId int) (int, error) {
	stmt := `SELECT COUNT(*) FROM comment_reactions
	WHERE comment_id = ? AND is_like = 0`
	var dislikes int
	row := c.DB.QueryRow(stmt, commentId)
	err := row.Scan(&dislikes)
	if err != nil {
		return 0, err
	}

	return dislikes, nil
}
