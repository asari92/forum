package repository

import (
	"database/sql"
	"errors"

	"forum/internal/entities"
)

const (
	activeStatus  = "read"
	passiveStatus = "unread"
)

type PostReactionSqlite3 struct {
	DB *sql.DB
}

func NewPostReactionSqlite3(db *sql.DB) *PostReactionSqlite3 {
	return &PostReactionSqlite3{DB: db}
}

func (r *PostReactionSqlite3) UpdateNotificationStatus(userID int) error {
	stmt := `UPDATE notifications
	SET status = ?
	WHERE user_id = ?`
	_, err := r.DB.Exec(stmt, activeStatus, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostReactionSqlite3) GetNotifications(userID int) ([]*entities.Notification, error) {
	stmt := `
	SELECT n.id,n.post_id,n.action_type, n.trigger_user_id,n.created, u.username, p.title, p.content FROM notifications as n 
	JOIN users as u ON n.trigger_user_id = u.id JOIN posts as p ON n.post_id = p.id
	WHERE n.user_id = ? and n.status = ?
	ORDER BY n.created DESC
	`

	rows, err := r.DB.Query(stmt, userID, passiveStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entities.ErrNoRecord
		} else {
			return nil, err
		}
	}
	notifications := []*entities.Notification{}
	for rows.Next() {
		notification := &entities.Notification{}
		err := rows.Scan(
			&notification.ID,
			&notification.PostID,
			&notification.Action,
			&notification.TriggerUserID,
			&notification.Created,
			&notification.TriggerUserName,
			&notification.PostTitle,
			&notification.PostContent)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *PostReactionSqlite3) AddNotification(userID, postID, triggerUserID int, actionType string) error {
	stmt := `INSERT INTO notifications (user_id, post_id, action_type, trigger_user_id, created)
	VALUES (?,?,?,?, datetime('now'))`
	_, err := r.DB.Exec(stmt, userID, postID, actionType, triggerUserID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostReactionSqlite3) RemoveNotification(userID, postID, triggerUserID int, actionType string) error {
	stmt := `DELETE FROM notifications
	WHERE user_id = ? AND post_id = ? AND trigger_user_id = ? AND action_type = ?`
	_, err := r.DB.Exec(stmt, userID, postID, triggerUserID, actionType)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostReactionSqlite3) UpdateNotification(userID, postID, triggerUserID int, oldAction, newAction string) error {
	stmt := `UPDATE notifications
	SET action_type = ?
	WHERE user_id = ? AND post_id = ? AND trigger_user_id = ? AND action_type = ?`
	_, err := r.DB.Exec(stmt, newAction, userID, postID, triggerUserID, oldAction)
	if err != nil {
		return err
	}
	return nil
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
