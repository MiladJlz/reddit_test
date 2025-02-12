package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/miladjlz/reddit_test/types"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PostStore interface {
	InsertPost(post *types.Post) (*types.Post, error)
	HasUserVoted(userID uuid.UUID, postID uuid.UUID) (bool, error)
	GetPosts() ([]types.Post, error)
	GetPostsByVote() ([]types.Post, error)
	GetPostsByDate(orderBy string) ([]types.Post, error)
}
type DBPostStore struct {
	postgresClient *gorm.DB
	redisClient    *redis.Client
}

func NewDBPostStore(pdb *gorm.DB, rdb *redis.Client) (*DBPostStore, error) {
	return &DBPostStore{pdb, rdb}, nil
}

func (s *DBPostStore) InsertPost(post *types.Post) (*types.Post, error) {
	result := s.postgresClient.Table("posts").Create(&post)
	if result.Error != nil {
		return nil, result.Error
	}
	return post, nil
}
func (s *DBPostStore) HasUserVoted(userID uuid.UUID, postID uuid.UUID) (bool, error) {
	var vote types.Vote
	result := s.postgresClient.Where("user_id = ? AND post_id = ?", userID, postID).First(&vote)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}
func (s *DBPostStore) GetPosts() ([]types.Post, error) {
	var posts []types.Post
	result := s.postgresClient.Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}
func (s *DBPostStore) GetPostsByScoreFromPostgres() ([]types.Post, error) {
	var posts []types.Post
	result := s.postgresClient.Where("vote_count >= ?", 5).Order("vote_count DESC").Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}
func (s *DBPostStore) GetPostsByVote() ([]types.Post, error) {
	var posts []types.Post
	exists, err := s.redisClient.Exists(context.Background(), "high_voted_posts").Result()
	fmt.Println(exists)
	if err != nil {
		return nil, err
	}
	if exists == 0 {
		return s.GetPostsByScoreFromPostgres()
	}
	postsJSON, err := s.redisClient.Get(context.Background(), "high_voted_posts").Bytes()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(postsJSON, &posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *DBPostStore) GetPostsByDate(orderBy string) ([]types.Post, error) {
	var posts []types.Post
	var result *gorm.DB
	if orderBy == "ASC" {
		result = s.postgresClient.Table("posts").Order("created_at ASC").Find(&posts)
	} else {
		result = s.postgresClient.Table("posts").Order("created_at DESC").Find(&posts)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}
