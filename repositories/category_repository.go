package repositories

import (
	"database/sql"
	"errors"

	"forum/entities"
)

type CategoryRepository struct {
	DB *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (r *CategoryRepository) Insert(name string) (int, error) {
	stmt := `INSERT INTO categories (name) VALUES (?)`
	result, err := r.DB.Exec(stmt, name)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *CategoryRepository) Get(categoryId int) (*entities.Category, error) {
	stmt := `SELECT id, name FROM categories WHERE id = ?`
	row := r.DB.QueryRow(stmt, categoryId)
	category := &entities.Category{}
	err := row.Scan(&category.ID, &category.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return category, nil
}

func (r *CategoryRepository) GetAll() ([]*entities.Category, error) {
	stmt := `SELECT id, name FROM categories`
	rows, err := r.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*entities.Category{}
	for rows.Next() {
		c := &entities.Category{}
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) Delete(categoryId int) error {
	stmt := `DELETE FROM categories WHERE id = ?`
	_, err := r.DB.Exec(stmt, categoryId)
	return err
}

func (r *CategoryRepository) GetCategoriesForPost(postId int) ([]*entities.Category, error) {
	stmt := `SELECT c.id, c.name FROM categories c
	         INNER JOIN post_categories pc ON c.id = pc.category_id
	         WHERE pc.post_id = ?`
	rows, err := r.DB.Query(stmt, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entities.Category
	for rows.Next() {
		c := &entities.Category{}
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) DeleteCategoriesForPost(postId int) error {
	stmt := `DELETE FROM post_categories WHERE post_id = ?`
	_, err := r.DB.Exec(stmt, postId)
	return err
}
