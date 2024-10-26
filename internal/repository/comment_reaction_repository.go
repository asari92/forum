package repository

import (
	"database/sql"

	"forum/internal/entities"
)

type CommentReactionSqlite3 struct {
	DB *sql.DB
}

func NewCommentReactionSqlite3(db *sql.DB) *CommentReactionSqlite3 {
	return &CommentReactionSqlite3{DB: db}
}

//	прежде чем добавить реакцию я пробую достать из бд реакцию пользователя под этот коммент,
//
// если передаваемая через пост запрос реакция и та реакция которую я достал одинаковые
// то просто удаляю эту запись
func (r *CommentReactionSqlite3) AddReaction(userID, commentID int, isLike bool) error {
	stmt := `INSERT INTO comment_reactions (user_id, comment_id,is_like)
	VALUES (?,?,?)
	ON CONFLICT(user_id, comment_id) DO UPDATE SET is_like = ?`
	_, err := r.DB.Exec(stmt, userID, commentID, isLike, isLike)
	return err
}

func (r *CommentReactionSqlite3) RemoveReaction(userID, commentID int) error {
	stmt := `DELETE FROM comment_reactions
	WHERE user_id = ? AND comment_id = ?`
	_, err := r.DB.Exec(stmt, userID, commentID)
	if err != nil {
		return err
	}
	return nil
}

func (r *CommentReactionSqlite3) GetUserReaction(userID, commentID int) (*entities.CommentReaction, error) {
	stmt := `SELECT is_like FROM comment_reactions WHERE user_id = ? AND comment_id = ?`
	row := r.DB.QueryRow(stmt, userID, commentID)

	commentReaction := &entities.CommentReaction{}
	err := row.Scan(&commentReaction.IsLike)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return commentReaction, nil
}

func (r *CommentReactionSqlite3) GetLikesCount(commentID int) (int, error) {
	stmt := `SELECT COUNT(is_like) FROM comment_reactions
	WHERE comment_id = ? AND is_like = 1`
	row := r.DB.QueryRow(stmt, commentID)

	var likes int
	err := row.Scan(&likes)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return likes, nil
}

func (r *CommentReactionSqlite3) GetDislikesCount(commentID int) (int, error) {
	stmt := `SELECT COUNT(is_like) FROM comment_reactions
	WHERE comment_id = ? AND is_like = 0`
	var dislikes int
	row := r.DB.QueryRow(stmt, commentID)
	err := row.Scan(&dislikes)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return dislikes, nil
}
