package service

import (
	"forum/internal/entities"
	"forum/internal/repository"
)

type User interface {
	NewUserAuthForm() userAuthForm
	NewAccountPasswordUpdateForm() accountPasswordUpdateForm
	Insert(form *userAuthForm) error
	Authenticate(form *userAuthForm) (int, error)
	UserExists(id int) (bool, error)
	GetUserByID(id int) (*entities.User, error)
	UpdatePassword(userID int, form *accountPasswordUpdateForm) error
}

type Post interface {
	NewPostCreateForm() postCreateForm
	GetPostDTO(postID int, userID int) (*PostDTO, error)
	GetAllPaginatedPostsDTO(page, pageSize int, paginationURL string) (*PostsDTO, error)
	GetUserPostsDTO(userID, page, pageSize int, paginationURL string) (*PostsDTO, error)
	GetUserLikedPostsDTO(userID, page, pageSize int, paginationURL string) (*PostsDTO, error)
	GetFilteredPaginatedPostsDTO(form *postCreateForm, page, pageSize int, paginationURL string) (*PostsDTO, error)
	CreatePostWithCategories(form *postCreateForm, userID int) (int, []*entities.Category, error)
	DeletePost(postID int) error
}

type Reaction interface {
	NewReactionForm() reactionForm
	UpdatePostReaction(userID, postID int, form *reactionForm) error
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
