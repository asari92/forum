package service

import (
	"errors"
	"fmt"
	"strings"

	"forum/internal/entities"
	"forum/internal/repository"
	"forum/pkg/validator"
)

// Use Case структура
type PostUseCase struct {
	categoryRepo        repository.CategoryRepository
	commentRepo         repository.CommentRepository
	commentReactionRepo repository.CommentReactionRepository
	postRepo            repository.PostRepository
	postReactionRepo    repository.PostReactionRepository
	userRepo            repository.UserRepository
}

type PostDTO struct {
	Post         *entities.Post
	Categories   []*entities.Category
	Likes        int
	Dislikes     int
	Comments     []*entities.Comment
	UserReaction *entities.PostReaction
}

type PostsDTO struct {
	User          *entities.User
	Posts         []*entities.Post
	HasNextPage   bool
	CurrentPage   int
	PaginationURL string
	Categories    []*entities.Category
	Header        string
}

type postCreateForm struct {
	Title      string
	Content    string
	Categories []int
	validator.Validator
}

func NewPostUseCase(repo *repository.Repository) *PostUseCase {
	return &PostUseCase{
		categoryRepo:        repo.CategoryRepository,
		commentRepo:         repo.CommentRepository,
		commentReactionRepo: repo.CommentReactionRepository,
		postRepo:            repo.PostRepository,
		postReactionRepo:    repo.PostReactionRepository,
		userRepo:            repo.UserRepository,
	}
}

func (uc *PostUseCase) NewPostCreateForm() postCreateForm {
	return postCreateForm{Categories: []int{}}
}

func (uc *PostUseCase) GetPostDTO(postID int, userID int) (*PostDTO, error) {
	post, err := uc.postRepo.GetPost(postID)
	if err != nil {
		return nil, err
	}

	categories, err := uc.categoryRepo.GetCategoriesForPost(postID)
	if err != nil {
		return nil, err
	}

	likes, dislikes, err := uc.postReactionRepo.GetReactionsCount(postID)
	if err != nil {
		return nil, err
	}

	var userReaction *entities.PostReaction
	if userID > 0 {
		userReaction, err = uc.postReactionRepo.GetUserReaction(userID, postID) // Получите реакцию пользователя
		if err != nil {
			return nil, err
		}
	}

	comments, err := uc.commentRepo.GetComments(postID)
	if err != nil {
		return nil, err
	}

	for _, comment := range comments {
		if userID > 0 {
			userReaction, err := uc.commentReactionRepo.GetUserReaction(userID, comment.ID)
			if err != nil {
				return nil, err
			}
			if userReaction != nil {
				if userReaction.IsLike {
					comment.UserReaction = 1
				} else {
					comment.UserReaction = -1
				}
			}
		}

		like, err := uc.commentReactionRepo.GetLikesCount(comment.ID)
		if err != nil {
			return nil, err
		}
		comment.Like = like

		dislike, err := uc.commentReactionRepo.GetDislikesCount(comment.ID)
		if err != nil {
			return nil, err
		}
		comment.Dislike = dislike
	}

	return &PostDTO{
		Post:         post,
		Categories:   categories,
		Likes:        likes,
		Dislikes:     dislikes,
		Comments:     comments,
		UserReaction: userReaction,
	}, nil
}

// Получение постов пользователя с пагинацией
func (uc *PostUseCase) GetUserPostsDTO(userID, page, pageSize int, paginationURL string) (*PostsDTO, error) {
	posts, err := uc.postRepo.GetUserPaginatedPosts(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepo.Get(userID)
	if err != nil {
		if errors.Is(err, entities.ErrNoRecord) {
			if len(posts) == 0 {
				return nil, err
			} else {
				user = &entities.User{Username: "Deleted User", ID: userID}
			}
		} else {
			return nil, err
		}
	}

	hasNextPage := len(posts) > pageSize
	if hasNextPage {
		posts = posts[:pageSize]
	}

	return &PostsDTO{
		User:          user,
		Posts:         posts,
		HasNextPage:   hasNextPage,
		CurrentPage:   page,
		PaginationURL: paginationURL,
	}, nil
}

// Получение постов, которые пользователь лайкнул
func (uc *PostUseCase) GetUserLikedPostsDTO(userID, page, pageSize int, paginationURL string) (*PostsDTO, error) {
	// Получаем посты, которые пользователь лайкнул
	posts, err := uc.postRepo.GetUserLikedPaginatedPosts(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	// Получаем пользователя
	user, err := uc.userRepo.Get(userID)
	if err != nil {
		return nil, err
	}

	// Проверяем наличие следующей страницы
	hasNextPage := len(posts) > pageSize
	if hasNextPage {
		posts = posts[:pageSize]
	}

	return &PostsDTO{
		User:          user,
		Posts:         posts,
		HasNextPage:   hasNextPage,
		CurrentPage:   page,
		PaginationURL: paginationURL,
	}, nil
}

// Получение всех постов с пагинацией
func (uc *PostUseCase) GetAllPaginatedPostsDTO(page, pageSize int, paginationURL string) (*PostsDTO, error) {
	// Получаем посты для нужной страницы.
	posts, err := uc.postRepo.GetAllPaginatedPosts(page, pageSize)
	if err != nil {
		return nil, err
	}

	// Проверяем, есть ли следующая страница
	hasNextPage := len(posts) > pageSize
	// Если постов больше, чем pageSize, обрезаем список до pageSize
	if hasNextPage {
		posts = posts[:pageSize]
	}

	categories, err := uc.categoryRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return &PostsDTO{
		Posts:         posts,
		HasNextPage:   hasNextPage,
		CurrentPage:   page,
		PaginationURL: paginationURL,
		Categories:    categories,
	}, nil
}

func (uc *PostUseCase) GetFilteredPaginatedPostsDTO(form *postCreateForm, page, pageSize int, paginationURL string) (*PostsDTO, error) {
	allCategories, err := uc.categoryRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Проверка категорий
	form.validateCategories(allCategories)

	PostsDTO := &PostsDTO{
		Header:        "All posts",
		CurrentPage:   page,
		PaginationURL: paginationURL,
		Categories:    allCategories,
	}

	if !form.Valid() {
		posts, err := uc.postRepo.GetAllPaginatedPosts(page, pageSize)
		if err != nil {
			return nil, err
		}

		// Проверяем, есть ли следующая страница
		hasNextPage := len(posts) > pageSize
		// Если постов больше, чем pageSize, обрезаем список до pageSize
		if hasNextPage {
			posts = posts[:pageSize]
		}

		PostsDTO.Posts = posts
		PostsDTO.HasNextPage = hasNextPage

		return PostsDTO, entities.ErrInvalidCredentials
	}

	// Получаем посты с пагинацией
	posts, err := uc.postRepo.GetPaginatedPostsByCategory(form.Categories, page, pageSize)
	if err != nil {
		return nil, err
	}

	// Проверяем, есть ли следующая страница
	hasNextPage := len(posts) > pageSize
	// Если постов больше, чем pageSize, обрезаем список до pageSize
	if hasNextPage {
		posts = posts[:pageSize]
	}

	// Собираем названия категорий, совпадающие с выбранными categoryIDs
	var filteredCategoryNames []string
	for _, category := range allCategories {
		for _, selectedID := range form.Categories {
			if category.ID == selectedID {
				filteredCategoryNames = append(filteredCategoryNames, category.Name)
			}
		}
	}

	// Если категории выбраны, создаем строку заголовка с их названиями
	if len(filteredCategoryNames) > 0 {
		PostsDTO.Header = fmt.Sprintf("Posts filtered by: %s", strings.Join(filteredCategoryNames, ", "))
	}

	PostsDTO.Posts = posts
	PostsDTO.HasNextPage = hasNextPage

	return PostsDTO, nil
}

// Создание поста с категориями
func (uc *PostUseCase) CreatePostWithCategories(form *postCreateForm, userID int) (int, []*entities.Category, error) {
	// валидировать все данные
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")

	allCategories, err := uc.categoryRepo.GetAll()
	if err != nil {
		return 0, nil, err
	}

	form.validateCategories(allCategories)

	if !form.Valid() {
		return 0, allCategories, entities.ErrInvalidCredentials
	}

	postID, err := uc.postRepo.InsertPostWithCategories(form.Title, form.Content, userID, form.Categories)
	if err != nil {
		return 0, allCategories, err
	}

	return postID, allCategories, nil
}

// Удаление поста
func (uc *PostUseCase) DeletePost(postID int) error {
	return uc.postRepo.DeletePost(postID)
}

func (form *postCreateForm) validateCategories(allCategories []*entities.Category) {
	categoryIDs := []int{}
	if len(form.Categories) == 0 {
		form.AddFieldError("categories", "Need one or more category")
	} else {
		categoryMap := make(map[int]bool)
		for _, category := range allCategories {
			categoryMap[category.ID] = true
		}

		for _, c := range form.Categories {
			if _, exists := categoryMap[c]; exists {
				categoryIDs = append(categoryIDs, c)
			} else {
				form.AddFieldError("categories", "One or more categories are invalid")
			}
		}
	}

	if len(categoryIDs) > 0 {
		form.Categories = categoryIDs
	} else {
		form.AddFieldError("categories", "Need one or more category")
	}
}
