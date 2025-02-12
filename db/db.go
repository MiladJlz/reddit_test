package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/miladjlz/reddit_test/types"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type DBConfig struct {
	PostgresClient *gorm.DB
	RedisClient    *redis.Client
	PostStore      *DBPostStore
	UserStore      *PostgresUserStore
	VoteStore      *PostgresVoteStore
	ServerPort     string
}

func NewDBConfig() (*DBConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, errors.New("failed to initializing postgres")
	}

	port := os.Getenv("SERVER_PORT")
	postgresDB, err := initPostgres()
	if err != nil {
		return nil, errors.New("failed to initializing postgres")
	}
	redisDB, err := initRedis()
	if err != nil {
		return nil, errors.New("failed to initializing redis")
	}
	postStore, err := NewDBPostStore(postgresDB, redisDB)
	if err != nil {
		return nil, errors.New("failed to create post store")
	}
	userStore, err := NewPostgresUserStore(postgresDB)
	if err != nil {
		return nil, errors.New("failed to create user store")
	}
	voteStore, err := NewPostgresVoteStore(postgresDB)
	if err != nil {
		return nil, errors.New("failed to create vote store")
	}
	return &DBConfig{
		PostgresClient: postgresDB,
		RedisClient:    redisDB,
		PostStore:      postStore,
		UserStore:      userStore,
		VoteStore:      voteStore,
		ServerPort:     port,
	}, nil
}
func initPostgres() (*gorm.DB, error) {
	postgresHost := os.Getenv("POSTGRES_DB_HOST")
	postgresPort := os.Getenv("POSTGRES_DB_PORT")
	postgresUser := os.Getenv("POSTGRES_DB_USER")
	postgresPass := os.Getenv("POSTGRES_DB_PASS")
	postgresDBName := os.Getenv("POSTGRES_DB_NAME")
	sslMode := os.Getenv("SSLMODE")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		postgresHost, postgresUser, postgresPass, postgresDBName, postgresPort, sslMode)

	pdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to connect to postgres")
	}
	if err := pdb.AutoMigrate(&types.User{}, &types.Post{}, &types.Vote{}); err != nil {
		return nil, errors.New("failed to migrate database")
	}
	return pdb, nil
}
func initRedis() (*redis.Client, error) {
	redisHost := os.Getenv("REDIS_DB_HOST")
	redisPort := os.Getenv("REDIS_DB_PORT")
	redisPassword := os.Getenv("REDIS_DB_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, errors.New("failed to connect to Redis")
	}
	return rdb, nil
}
