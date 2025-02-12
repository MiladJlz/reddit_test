package db

import (
	"github.com/miladjlz/reddit_test/types"
	"gorm.io/gorm"
)

type VoteStore interface {
	InsertVote(vote *types.Vote) (*types.Vote, error)
}
type PostgresVoteStore struct {
	client *gorm.DB
}

func NewPostgresVoteStore(db *gorm.DB) (*PostgresVoteStore, error) {
	return &PostgresVoteStore{db}, nil
}

func (s *PostgresVoteStore) InsertVote(vote *types.Vote) (*types.Vote, error) {
	tx := s.client.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	result := tx.Table("votes").Create(&vote)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	if result := tx.Table("posts").Where("id = ?", vote.PostID).Update("vote_count", gorm.Expr("vote_count + ?", 1)); result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	return vote, tx.Commit().Error
}
