package service

import (
	"errors"
	"forum/entities"
	"forum/repository"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")

	ErrInvalidCredentials = errors.New("models: invalid credentials")

	ErrDuplicateEmail = errors.New("models: duplicate email")
)

type User interface {
	Insert(username, email, password string) error
	Authenticate(email, password string) (int, error)
	UserExists(id int) (bool, error)
	GetUserByID(id int) (*entities.User, error)
	UpdatePassword(id int, currentPassword, newPassword string) error
}

type Post interface {
	CreatePostWithCategories(title, content string, userID int, categoryIDs []int) (int, error)
	GetPost(postID int) (*entities.Post, error)
	GetPaginatedPostsByCategory(categoryIDs []int, page, pageSize int) ([]*entities.Post, error)
	GetUserPaginatedPosts(userID, page, pageSize int) ([]*entities.Post, error)
	GetUserLikedPaginatedPosts(userID, page, pageSize int) ([]*entities.Post, error)
	GetAllPaginatedPosts(page, pageSize int) ([]*entities.Post, error)
	DeletePost(postID int) error
}

type PostReaction interface {
	AddReaction(userID, postID int, isLike bool) error
	RemoveReaction(userID, postID int) error
	GetUserReaction(userID, postID int) (*entities.PostReaction, error)
	GetReactionsCount(postID int) (int, int, error)
}

type Comment interface {
	AddComment(postID, userID int, content string) error
	GetPostComments(postID int) ([]*entities.Comment, error)
}

type CommentReaction interface {
	AddReaction(userID, commentID int, isLike bool) error
	RemoveReaction(userID, postID int) error
	GetUserReaction(userID, commentID int) (*entities.CommentReaction, error)
	GetLikesCount(commentID int) (int, error)
	GetDislikesCount(commentID int) (int, error)
}

type Category interface {
	Insert(name string) (int, error)
	Get(categoryId int) (*entities.Category, error)
	GetAll() ([]*entities.Category, error)
	Delete(categoryId int) error
	GetCategoriesForPost(postId int) ([]*entities.Category, error)
	DeleteCategoriesForPost(postId int) error
}

type Service struct {
	User
	Post
	PostReaction
	Comment
	CommentReaction
	Category
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:            NewUserUseCase(repos.UserRepository),
		Post:            NewPostUseCase(repos.PostRepository),
		PostReaction:    NewPostReactionUseCase(repos.PostReactionRepository),
		Comment:         NewCommentUseCase(repos.CommentRepository),
		CommentReaction: NewCommentReactionUseCase(repos.CommentReactionRepository),
		Category:        NewCategoryUseCase(repos.CategoryRepository),
	}
}