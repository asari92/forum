package repository

import (
	"database/sql"

	"forum/entities"
)

type PostReactionSqlite3 struct {
	DB *sql.DB
}

func NewPostReactionSqlite3(db *sql.DB) *PostReactionSqlite3 {
	return &PostReactionSqlite3{DB: db}
}

func (r *PostReactionSqlite3) AddReaction(userID, postID int, isLike bool) error {
	stmt := `INSERT INTO post_reactions (user_id, post_id, is_like)
              VALUES (?, ?, ?)
              ON CONFLICT(user_id, post_id) DO UPDATE SET is_like = ?`
	_, err := r.DB.Exec(stmt, userID, postID, isLike, isLike)
	return err
}

func (r *PostReactionSqlite3) RemoveReaction(userID, postID int) error {
	stmt := `DELETE FROM post_reactions WHERE user_id = ? AND post_id = ?`
	_, err := r.DB.Exec(stmt, userID, postID)
	return err
}

func (r *PostReactionSqlite3) GetUserReaction(userID, postID int) (*entities.PostReaction, error) {
	stmt := `SELECT is_like FROM post_reactions WHERE post_id = ? AND user_id = ?`
	row := r.DB.QueryRow(stmt, postID, userID)

	var reaction entities.PostReaction
	err := row.Scan(&reaction.IsLike)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Если реакции нет, возвращаем nil
		}
		return nil, err
	}

	return &reaction, nil
}

func (r *PostReactionSqlite3) GetReactionsCount(postID int) (likes int, dislikes int, err error) {
	stmt := `SELECT COALESCE(SUM(CASE WHEN is_like = 1 THEN 1 ELSE 0 END), 0) AS likes,
                     COALESCE(SUM(CASE WHEN is_like = 0 THEN 1 ELSE 0 END), 0) AS dislikes
              FROM post_reactions WHERE post_id = ?`
	row := r.DB.QueryRow(stmt, postID)
	err = row.Scan(&likes, &dislikes)
	return
}
