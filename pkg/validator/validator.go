package validator

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	maxImageSize = 20 * 1024 * 1024
)

var (
	UsernameRX = regexp.MustCompile(`^[^._ ](?:[\w-]|\.[\w-])+[^._ ]$`)
	EmailRX    = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	PasswordRX = regexp.MustCompile("[0-9a-zA-Z!_.@#$%^&*]{8,}")
	TextRX     = regexp.MustCompile(`^[а-яА-ЯёЁa-zA-Z0-9.,:;!?'"()\-–—\[\]{}<>/|@#$%^&*+=_~\s]+$`)
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) AddFieldError(key, message string) {
	// Note: We need to initialize the map first, if it isn't already
	// initialized.
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func NotHaveAnySpaces(value string) bool {
	return strings.Join(strings.Fields(value), "") == value
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// ValidateID выполняет полную проверку ID
func ValidateID(id string) (int, error) {
	// Проверяем, начинается ли строка с "+" или "-"
	if id == "" {
		return 0, errors.New("ID не может быть пустым")
	}
	if id[0] == '+' || id[0] == '-' {
		return 0, errors.New("ID не может содержать знак '+' или '-'")
	}

	// Проверяем на ведущие нули
	if len(id) > 1 && id[0] == '0' {
		return 0, errors.New("ID содержит незначащие нули")
	}

	// Проверяем, что ID состоит только из цифр
	for _, r := range id {
		if !unicode.IsDigit(r) {
			return 0, errors.New("ID может содержать только цифры")
		}
	}

	// Преобразуем строку в число
	ID, err := strconv.Atoi(id)
	if err != nil || ID < 1 {
		return 0, errors.New("ID должен быть положительным числом")
	}

	return ID, nil
}

func ValidateImageFiles(files []*multipart.FileHeader) error {
	for _, fileHeader := range files {
		if fileHeader.Size > maxImageSize {
			return errors.New("file size larger than max")
		}

		file, err := fileHeader.Open()
		if err != nil {
			return fmt.Errorf("error opening file %s: %v", fileHeader.Filename, err)
		}
		defer file.Close()

		buf := make([]byte, 512)
		_, err = file.Read(buf)
		if err != nil {
			return fmt.Errorf("error reading file %s: %v", fileHeader.Filename, err)
		}

		mimeType := http.DetectContentType(buf)
		validMimeTypes := map[string]bool{
			"image/jpeg": true,
			"image/jpg":  true,
			"image/png":  true,
			"image/gif":  true,
		}

		if !validMimeTypes[mimeType] {
			return errors.New("unsupported file type")
		}

	}
	return nil
}
