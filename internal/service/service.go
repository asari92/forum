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
	GetPaginatedPostsByCategory(categoryIDs []int, page, pageSize int) ([]*entities.Post, error)
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
	RemoveReaction(userID, commentID int) error
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
		Post:            NewPostUseCase(repos),
		PostReaction:    NewPostReactionUseCase(repos.PostReactionRepository),
		Comment:         NewCommentUseCase(repos.CommentRepository),
		CommentReaction: NewCommentReactionUseCase(repos.CommentReactionRepository),
		Category:        NewCategoryUseCase(repos.CategoryRepository),
	}
}
