package models

import "database/sql"

type PostReactionModelInterface interface {
	AddReaction(userID, postID int, isLike bool) error
	RemoveReaction(userID, postID int) error
	GetUserReaction(userID, postID int) (*PostReaction, error)
	GetReactionsCount(postID int) (likes int, dislikes int, err error)
}

type PostReaction struct {
	IsLike bool
}

type PostReactionModel struct {
	DB *sql.DB
}

// Добавление или обновление реакции
func (m *PostReactionModel) AddReaction(userID, postID int, isLike bool) error {
	stmt := `INSERT INTO post_reactions (user_id, post_id, is_like)
              VALUES (?, ?, ?)
              ON CONFLICT(user_id, post_id) DO UPDATE SET is_like = ?`
	_, err := m.DB.Exec(stmt, userID, postID, isLike, isLike)
	return err
}

// Удаление реакции
func (m *PostReactionModel) RemoveReaction(userID, postID int) error {
	stmt := `DELETE FROM post_reactions WHERE user_id = ? AND post_id = ?`
	_, err := m.DB.Exec(stmt, userID, postID)
	return err
}

func (m *PostReactionModel) GetUserReaction(userID, postID int) (*PostReaction, error) {
	stmt := `SELECT is_like FROM post_reactions WHERE post_id = ? AND user_id = ?`
	row := m.DB.QueryRow(stmt, postID, userID)

	var reaction PostReaction
	err := row.Scan(&reaction.IsLike)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Если реакции нет, возвращаем nil
		}
		return nil, err
	}

	return &reaction, nil
}

// Получение количества лайков и дизлайков для поста
func (m *PostReactionModel) GetReactionsCount(postID int) (likes int, dislikes int, err error) {
	stmt := `SELECT COALESCE(SUM(CASE WHEN is_like = 1 THEN 1 ELSE 0 END), 0) AS likes,
                     COALESCE(SUM(CASE WHEN is_like = 0 THEN 1 ELSE 0 END), 0) AS dislikes
              FROM post_reactions WHERE post_id = ?`
	row := m.DB.QueryRow(stmt, postID)
	err = row.Scan(&likes, &dislikes)
	return
}
 