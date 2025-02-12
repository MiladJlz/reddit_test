package db

import (
	"github.com/miladjlz/reddit_test/types"
	"gorm.io/gorm"
)

type UserStore interface {
	InsertUser(user *types.User) (*types.User, error)
}
type PostgresUserStore struct {
	client *gorm.DB
}

func NewPostgresUserStore(db *gorm.DB) (*PostgresUserStore, error) {
	return &PostgresUserStore{db}, nil
}

func (s *PostgresUserStore) InsertUser(User *types.User) (*types.User, error) {
	result := s.client.Table("users").Create(&User)
	if result.Error != nil {
		return nil, result.Error
	}
	return User, nil
}
