package service

import (
	"testing"

	"forum/internal/repository"
	"forum/pkg/assert"
)

func TestUserModelExists(t *testing.T) {
	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			m := NewUserUseCase(repository.NewUserSqlite3(db))
			exists, err := m.UserExists(tt.userID)
			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}