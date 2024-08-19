package models

import (
	"database/sql"
)

const DefaultCategory = 1

type Category struct {
	ID   int
	Name string
}

type CategoryModel struct {
	DB *sql.DB
}

func (m *CategoryModel) Insert(name string) (int, error) {
	stmt := `INSERT INTO categories (name)
    VALUES(?)`

	result, err := m.DB.Exec(stmt, name)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *CategoryModel) GetAll() ([]*Category, error) {
	stmt := `SELECT id, name FROM categories`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*Category{}

	for rows.Next() {
		c := &Category{}
		err = rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}

		categories = append(categories, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (m *CategoryModel) Delete(categoryId int) error {
	stmt := `DELETE FROM categories WHERE id = ?`

	_, err := m.DB.Exec(stmt, categoryId)
	if err != nil {
		return err
	}

	return nil
}

// func (m *CategoryModel) InsertCategoriesForPost(postId int, categoryIDs []int) error {
// 	stmt := `INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`

// 	// Начинаем транзакцию для выполнения нескольких вставок
// 	tx, err := m.DB.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	// Подготовленный запрос для вставки данных
// 	stmtInsert, err := tx.Prepare(stmt)
// 	if err != nil {
// 		tx.Rollback() // откат транзакции в случае ошибки
// 		return err
// 	}
// 	defer stmtInsert.Close()

// 	// Вставляем каждую категорию для поста
// 	for _, categoryID := range categoryIDs {
// 		_, err := stmtInsert.Exec(postId, categoryID)
// 		if err != nil {
// 			tx.Rollback() // откат транзакции в случае ошибки
// 			return err
// 		}
// 	}

// 	// Завершаем транзакцию
// 	err = tx.Commit()
// 	if err != nil {
// 		tx.Rollback() // откат транзакции в случае ошибки
// 		return err
// 	}

// 	return nil
// }

func (m *CategoryModel) GetCategoriesForPost(postId int) ([]*Category, error) {
	stmt := `SELECT c.id, c.name 
	FROM categories c
	INNER JOIN post_categories pc ON c.id = pc.category_id
    WHERE pc.post_id = ?`

	rows, err := m.DB.Query(stmt, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*Category{}

	for rows.Next() {
		c := &Category{}
		err = rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}

		categories = append(categories, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (m *CategoryModel) DeleteCategoriesForPost(postId int) error {
	stmt := `DELETE FROM post_categories WHERE post_id = ?`

	_, err := m.DB.Exec(stmt, postId)
	if err != nil {
		return err
	}

	return nil
}
