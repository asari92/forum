package service

import (
	"mime/multipart"

	"forum/internal/entities"
	"forum/internal/repository"
)

type User interface {
	NewUserAuthForm() userAuthForm
	NewAccountPasswordUpdateForm() accountPasswordUpdateForm
	Insert(username, email, password, role string) (int, error)
	Authenticate(email, password string) (*entities.User, error)
	OauthAuthenticate(email string) (*entities.User, error)
	UserExists(id int) (bool, error)
	GetUserByID(id int) (*entities.User, error)
	UpdatePassword(userID int, form *accountPasswordUpdateForm) error
}

type Post interface {
	NewPostCreateForm() postCreateForm
	NewCommentForm() CommentForm
	GetPostDTO(postID int, userID int) (*PostDTO, error)
	GetCommentedPostDTO(postID int, userID int) (*PostDTO, error)
	GetAllPaginatedPostsDTO(page, pageSize int, paginationURL string) (*PostsDTO, error)
	GetAllPaginatedUnapprovedPostsDTO(page, pageSize int, paginationURL string) (*PostsDTO, error)
	GetUserPostsDTO(userID, page, pageSize int, paginationURL string) (*PostsDTO, error)
	GetUserCommentedPostsDTO(userID, page, pageSize int, paginationURL string) (*PostsDTO, error)
	GetUserLikedPostsDTO(userID, page, pageSize int, paginationURL string) (*PostsDTO, error)
	GetFilteredPaginatedPostsDTO(form *postCreateForm, page, pageSize int, paginationURL string) (*PostsDTO, error)
	CreatePostWithCategories(form *postCreateForm, files []*multipart.FileHeader, userID int) (int, []*entities.Category, error)
	UpdatePostWithImage(form *postCreateForm, postID int, files []*multipart.FileHeader, userID int) error
	DeleteComment(commentID, userID int) error
	UpdateComment(form *CommentForm, commentID, userID int) error
	GetUserNotifications(userID int) ([]*entities.Notification, error)
	ApprovePost(postID int) error
	DeletePost(postID, userID int) error
}

type Reaction interface {
	NewReactionForm() reactionForm
	UpdatePostReaction(userID, postID int, form *reactionForm) error
	CreateReport(userID, postID int, reason string) error
	GetPostReport(postID int) (*entities.Report, error)
	GetAllPaginatedPostReportsDTO(page, pageSize int, paginationURL string) (*ReportsDTO, error)
}

type Category interface {
	Insert(name string) (int, error)
	Get(categoryId int) (*entities.Category, error)
	GetAll() ([]*entities.Category, error)
	Delete(categoryId int) error
}

type Service struct {
	User
	Post
	Reaction
	Category
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:     NewUserUseCase(repos.UserRepository),
		Post:     NewPostUseCase(repos),
		Reaction: NewReactionUseCase(repos),
		Category: NewCategoryUseCase(repos.CategoryRepository),
	}
}
