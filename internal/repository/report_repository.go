package repository

import (
	"database/sql"
	"errors"
	"time"

	"forum/internal/entities"
)

type ReportSqlite3 struct {
	DB *sql.DB
}

func NewReportSqlite3(db *sql.DB) *ReportSqlite3 {
	return &ReportSqlite3{
		DB: db,
	}
}

func (r *ReportSqlite3) CreateReport(userID, postID int, reason string) error {
	stmt := `INSERT INTO reports (user_id, post_id, reason, created)
			 VALUES (?, ?, ?, datetime('now'))`
	_, err := r.DB.Exec(stmt, userID, postID, reason)
	return err
}

func (r *ReportSqlite3) GetPostReport(postID int) (*entities.Report, error) {
	stmt := `SELECT id, user_id, post_id, reason, created  
	FROM reports
	WHERE post_id = ?`

	row := r.DB.QueryRow(stmt, postID)

	report := &entities.Report{}
	var created string

	err := row.Scan(&report.ID, &report.UserID, &report.PostID, &report.Reason, &created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entities.ErrNoRecord
		} else {
			return nil, err
		}
	}

	userTime, err := time.Parse("2006-01-02 15:04:05", created)
	if err != nil {
		return nil, err
	}
	report.Created = userTime.Format(time.RFC3339)

	return report, nil
}

func (r *ReportSqlite3) GetAllPaginatedPostReports(page, pageSize int) ([]*entities.Report, error) {
	offset := (page - 1) * pageSize // Вычисляем смещение для текущей страницы

	stmt := `SELECT reports.id, post_id, user_id, username, reason, reports.created FROM reports
	JOIN users ON reports.user_id = users.id
             LIMIT ? OFFSET ?`

	rows, err := r.DB.Query(stmt, pageSize+1, offset) // Лимит на одну запись больше
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reports := []*entities.Report{}
	var created string

	for rows.Next() {
		report := &entities.Report{}
		err = rows.Scan(&report.ID, &report.PostID, &report.UserID, &report.ReporterName, &report.Reason, &created)
		if err != nil {
			return nil, err
		}

		reportTime, err := time.Parse("2006-01-02 15:04:05", created)
		if err != nil {
			return nil, err
		}
		report.Created = reportTime.Format(time.RFC3339)

		reports = append(reports, report)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *ReportSqlite3) DeleteReport(userID, postID int) error {
	stmt := `DELETE FROM reports
			WHERE user_id = ? AND post_id = ?`
	_, err := r.DB.Exec(stmt, userID, postID)
	return err
}
