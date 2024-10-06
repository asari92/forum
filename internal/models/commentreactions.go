package models

import "database/sql"

type CommentReaction struct {
	IsLike bool
}

type CommentReactionModelInterface interface {
	AddReaction(userID, commentID int, isLike bool) error
	RemoveReaction(userID, commentID int) error
	GetLikesCount(commentID int) (int, error)
	GetDislikesCount(commentID int) (int, error)
	GetUserReaction(userID, commentID int) (*CommentReaction, error)
}

type CommentReactionModel struct {
	DB *sql.DB
}

func (c *CommentReactionModel) AddReaction(userID, commentID int, isLike bool) error {
	stmt := `INSERT INTO comment_reactions (user_id, comment_id,is_like)
	VALUES (?,?,?)
	ON CONFLICT(user_id, comment_id) DO UPDATE SET is_like = ?`
	_, err := c.DB.Exec(stmt, userID, commentID, isLike, isLike)
	if err != nil {
		return err
	}
	return nil
}

//  прежде чем добавить реакцию я пробую достать из бд реакцию пользователя под этот коммент
// если передаваемый через пост запрос реакцию и та реакция которую я достал одинаковые то просто удаляю эту запись

func (c *CommentReactionModel) RemoveReaction(userID, commentID int) error {
	stmt := `DELETE FROM comment_reactions
	WHERE user_id = ? AND comment_id = ?`
	_, err := c.DB.Exec(stmt, userID, commentID)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommentReactionModel) GetLikesCount(commentID int) (int, error) {
	stmt := `SELECT COUNT(is_like) FROM comment_reactions
	WHERE comment_id = ? AND is_like = 1`
	var likes int
	row := c.DB.QueryRow(stmt, commentID)
	err := row.Scan(&likes)
	if err != nil {
		return 0, err
	}
	return likes, nil
}

func (c *CommentReactionModel) GetDislikesCount(commentID int) (int, error) {
	stmt := `SELECT COUNT(is_like) FROM comment_reactions
	WHERE comment_id = ? AND is_like = 0`
	var dislikes int
	row := c.DB.QueryRow(stmt, commentID)
	err := row.Scan(&dislikes)
	if err != nil {
		return 0, err
	}

	return dislikes, nil
}

func (c *CommentReactionModel) GetUserReaction(userID, commentID int) (*CommentReaction, error) {

	stmt := `SELECT is_like FROM comment_reactions WHERE user_id = ? AND comment_id = ?`
	row := c.DB.QueryRow(stmt, userID, commentID)
	commentReaction := &CommentReaction{}
	err := row.Scan(&commentReaction.IsLike)
	if err != nil {
		return nil, err
	}
	return commentReaction, nil

}
