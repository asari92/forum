package service

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"forum/internal/entities"
	"forum/internal/repository"
	"forum/pkg/validator"

	"github.com/gofrs/uuid"
)

const uploadDir = "uploads"

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
	Images       []*entities.Image
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

type CommentForm struct {
	Content string
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

func (uc *PostUseCase) NewCommentForm() CommentForm {
	return CommentForm{}
}

func (uc *PostUseCase) GetPostDTO(postID int, userID int) (*PostDTO, error) {
	exists, err := uc.postRepo.Exists(postID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, entities.ErrNoRecord
	}

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

	images, err := uc.postRepo.GetImagesByPost(postID)
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
		Images:       images,
		Comments:     comments,
		UserReaction: userReaction,
	}, nil
}

func (uc *PostUseCase) GetCommentedPostDTO(postID int, userID int) (*PostDTO, error) {
	exists, err := uc.postRepo.Exists(postID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, entities.ErrNoRecord
	}

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

	images, err := uc.postRepo.GetImagesByPost(postID)
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

	comments, err := uc.commentRepo.GetUserCommentsByPosts(postID, userID)
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
		Images:       images,
		Comments:     comments,
		UserReaction: userReaction,
	}, nil
}

// Получение постов пользователя с пагинацией
func (uc *PostUseCase) GetUserPostsDTO(userID, page, pageSize int, paginationURL string) (*PostsDTO, error) {
	exists, err := uc.userRepo.Exists(userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, entities.ErrNoRecord
	}

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

func (uc *PostUseCase) GetUserNotifications(userID int) ([]*entities.Notification, error) {
	exists, err := uc.userRepo.Exists(userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, entities.ErrNoRecord
	}

	notifications, err := uc.postReactionRepo.GetNotifications(userID)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

// Получение постов, которые пользователь лайкнул
func (uc *PostUseCase) GetUserLikedPostsDTO(userID, page, pageSize int, paginationURL string) (*PostsDTO, error) {
	exists, err := uc.userRepo.Exists(userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, entities.ErrNoRecord
	}

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

func (uc *PostUseCase) GetUserCommentedPostsDTO(userID, page, pageSize int, paginationURL string) (*PostsDTO, error) {
	exists, err := uc.userRepo.Exists(userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, entities.ErrNoRecord
	}

	// Получаем посты, которые пользователь лайкнул
	posts, err := uc.postRepo.GetUserCommentedPosts(userID, page, pageSize)
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

func (uc *PostUseCase) GetAllPaginatedUnapprovedPostsDTO(page, pageSize int, paginationURL string) (*PostsDTO, error) {
	// Получаем посты для нужной страницы.
	posts, err := uc.postRepo.GetAllPaginatedUnapprovedPosts(page, pageSize)
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

		return PostsDTO, nil
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
func (uc *PostUseCase) CreatePostWithCategories(form *postCreateForm, files []*multipart.FileHeader, userID int) (int, []*entities.Category, error) {
	// валидировать все данные
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")

	form.CheckField(validator.Matches(form.Title, validator.TextRX), "title", "This field must contain only english or russian letters")
	form.CheckField(validator.Matches(form.Content, validator.TextRX), "content", "This field must contain only english or russian letters")

	allCategories, err := uc.categoryRepo.GetAll()
	if err != nil {
		return 0, nil, err
	}

	form.validateCategories(allCategories)

	if len(files) != 0 {
		err = validator.ValidateImageFiles(files)
		if err != nil {
			if err.Error() == entities.ErrUnsupportedFileType.Error() {
				form.AddFieldError("image", "The project requires handling JPEG, PNG, GIF images")
			} else if err.Error() == entities.ErrFileSizeTooLarge.Error() {
				form.AddFieldError("image", "Maximum file size limit is 20 MB")
			} else {
				return 0, allCategories, err
			}
		}
	}

	if !form.Valid() {
		return 0, allCategories, entities.ErrInvalidCredentials
	}

	exists, err := uc.userRepo.Exists(userID)
	if err != nil {
		return 0, allCategories, err
	}
	if !exists {
		return 0, allCategories, entities.ErrNoRecord
	}

	filePaths := []string{}
	if len(files) != 0 {
		filePaths, err = uploadImages(files)
		if err != nil {
			return 0, allCategories, err
		}
	}

	postID, err := uc.postRepo.InsertPostWithCategories(form.Title, form.Content, userID, form.Categories, filePaths)
	if err != nil {
		for _, filePath := range filePaths {
			err := os.Remove(filePath)
			if err != nil {
				slog.Error("deleting file", "error", fmt.Sprintf("failed to delete file %s: %v", filePath, err))
			}
		}
		return 0, allCategories, err
	}
	return postID, allCategories, nil
}

func (uc *PostUseCase) UpdatePostWithImage(form *postCreateForm, postID int, files []*multipart.FileHeader, userID int) error {
	// валидировать все данные
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")

	form.CheckField(validator.Matches(form.Title, validator.TextRX), "title", "This field must contain only english or russian letters")
	form.CheckField(validator.Matches(form.Content, validator.TextRX), "content", "This field must contain only english or russian letters")

	if len(files) != 0 {
		err := validator.ValidateImageFiles(files)
		if err != nil {
			if err.Error() == entities.ErrUnsupportedFileType.Error() {
				form.AddFieldError("image", "The project requires handling JPEG, PNG, GIF images")
			} else if err.Error() == entities.ErrFileSizeTooLarge.Error() {
				form.AddFieldError("image", "Maximum file size limit is 20 MB")
			} else {
				return err
			}
		}
	}

	if !form.Valid() {
		return entities.ErrInvalidCredentials
	}

	exists, err := uc.userRepo.Exists(userID)
	if err != nil {
		return err
	}
	if !exists {
		return entities.ErrNoRecord
	}

	filePaths := []string{}
	if len(files) != 0 {
		filePaths, err = uploadImages(files)
		if err != nil {
			return err
		}
	}

	err = uc.postRepo.UpdatePostWithImage(form.Title, form.Content, postID, filePaths)
	if err != nil {
		for _, filePath := range filePaths {
			err := os.Remove(filePath)
			if err != nil {
				slog.Error("deleting file", "error", fmt.Sprintf("failed to delete file %s: %v", filePath, err))
			}
		}
		return err
	}
	return nil
}

func (uc *PostUseCase) UpdateComment(form *CommentForm, commentID, userID int) error {
	// валидировать все данные

	form.CheckField(validator.NotBlank(form.Content), "comment", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Content, validator.TextRX), "comment", "This field must contain only english or russian letters")

	if !form.Valid() {
		return entities.ErrInvalidCredentials
	}

	exists, err := uc.userRepo.Exists(userID)
	if err != nil {
		return err
	}
	if !exists {
		return entities.ErrNoRecord
	}

	err = uc.commentRepo.UpdateComment(commentID, form.Content)
	if err != nil {
		return err
	}

	err = uc.postReactionRepo.UpdateNotificationTime(commentID)
	if err != nil {
		return err
	}

	return nil
}

func uploadImages(files []*multipart.FileHeader) ([]string, error) {
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("error creating upload directory: %v", err)
		}
	}

	pathFiles := []string{}

	// Сохраняем файлы в папку
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("error opening file %s: %v", fileHeader.Filename, err)
		}
		defer file.Close()

		fileId, err := uuid.NewV4()
		if err != nil {
			return nil, fmt.Errorf("error generating UUID: %v", err)
		}
		fileExtension := path.Ext(fileHeader.Filename)
		safeFilename := fmt.Sprintf("%s%s", fileId.String(), fileExtension)

		filePath := path.Join(uploadDir, safeFilename)

		outFile, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("error creating file %s: %v", safeFilename, err)
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, file)
		if err != nil {
			return nil, fmt.Errorf("error saving file %s: %v", safeFilename, err)
		}

		pathFiles = append(pathFiles, filePath)

	}

	return pathFiles, nil
}

func (uc *PostUseCase) ApprovePost(postID int) error {
	exists, err := uc.postRepo.Exists(postID)
	if err != nil {
		return err
	}
	if !exists {
		return entities.ErrNoRecord
	}

	return uc.postRepo.ApprovePost(postID)
}

// Удаление поста
func (uc *PostUseCase) DeletePost(postID, userID int) error {
	exists, err := uc.postRepo.Exists(postID)
	if err != nil {
		return err
	}
	if !exists {
		return entities.ErrNoRecord
	}

	post, err := uc.postRepo.GetPost(postID)
	if err != nil {
		return err
	}

	if post.UserID == userID {
		filePaths, err := uc.postRepo.GetImagesByPost(postID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
			} else {
				return err
			}
		}

		if len(filePaths) != 0 {
			for _, filePath := range filePaths {
				err := os.Remove(filePath.UrlImage)
				if err != nil {
					return err
				}
			}
		}

		return uc.postRepo.DeletePost(postID)
	}
	return nil
}

func (uc *PostUseCase) DeleteComment(commentID, userID int) error {
	comment, err := uc.commentRepo.GetComment(commentID)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	if comment.UserID == userID {
		return uc.commentRepo.DeleteComment(commentID)
	}
	return nil
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
