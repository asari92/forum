package repository

import (
	"database/sql"

	"forum/internal/entities"
)

type UserRepository interface {
	Exists(id int) (bool, error)
	Insert(username, email, password string) (int, error)
	OauthAuthenticate(email string) (int, error)
	Authenticate(email, password string) (int, error)
	Get(id int) (*entities.User, error)
	UpdatePassword(id int, currentPassword, newPassword string) error
}

type PostRepository interface {
	Exists(id int) (bool, error)
	InsertPostWithCategories(title, content string, userID int, categoryIDs []int, filePaths []string) (int, error)
	GetPost(postID int) (*entities.Post, error)
	GetImagesByPost(postID int) ([]*entities.Image, error)
	GetPaginatedPostsByCategory(categoryIDs []int, page, pageSize int) ([]*entities.Post, error)
	GetUserCommentedPosts(userId, page, pageSize int) ([]*entities.Post, error)
	GetUserPaginatedPosts(userID, page, pageSize int) ([]*entities.Post, error)
	GetUserLikedPaginatedPosts(userID, page, pageSize int) ([]*entities.Post, error)
	GetAllPaginatedPosts(page, pageSize int) ([]*entities.Post, error)
	DeletePost(postID int) error
	UpdatePost(title, content string) error
}

type PostReactionRepository interface {
	AddReaction(userID, postID int, isLike bool) error
	RemoveReaction(userID, postID int) error
	GetUserReaction(userID, postID int) (*entities.PostReaction, error)
	GetReactionsCount(postID int) (likes int, dislikes int, err error)
}

type CommentRepository interface {
	Exists(id int) (bool, error)
	InsertComment(postID, userID int, content string) error
	GetComments(postID int) ([]*entities.Comment, error)
	GetUserCommentsByPosts(postId, userId int) ([]*entities.Comment, error)
}

type CommentReactionRepository interface {
	AddReaction(userID, commentID int, isLike bool) error
	RemoveReaction(userID, commentID int) error
	GetUserReaction(userID, commentID int) (*entities.CommentReaction, error)
	GetLikesCount(commentID int) (int, error)
	GetDislikesCount(commentID int) (int, error)
}

type CategoryRepository interface {
	Exists(id int) (bool, error)
	Insert(name string) (int, error)
	Get(categoryId int) (*entities.Category, error)
	GetAll() ([]*entities.Category, error)
	Delete(categoryId int) error
	GetCategoriesForPost(postId int) ([]*entities.Category, error)
	DeleteCategoriesForPost(postId int) error
}

type Repository struct {
	UserRepository
	PostRepository
	PostReactionRepository
	CommentRepository
	CommentReactionRepository
	CategoryRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		UserRepository:            NewUserSqlite3(db),
		PostRepository:            NewPostSqlite3(db),
		PostReactionRepository:    NewPostReactionSqlite3(db),
		CommentRepository:         NewCommentSqlite3(db),
		CommentReactionRepository: NewCommentReactionSqlite3(db),
		CategoryRepository:        NewCategorySqlite3(db),
	}
}
