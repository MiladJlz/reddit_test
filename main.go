package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/miladjlz/reddit_test/api"
	"github.com/miladjlz/reddit_test/db"
	"log"
	"sync"
	"time"
)

func main() {
	config, err := db.NewDBConfig()
	if err != nil {
		log.Fatal(err)
	}
	var (
		app         = fiber.New()
		postHandler = api.NewPostHandler(config.PostStore)
		voteHandler = api.NewVoteHandler(config.VoteStore, config.PostStore)
		userHandler = api.NewUserHandler(config.UserStore)
		v1          = app.Group("/v1")
	)

	//vote route
	v1.Post("/vote/:id", voteHandler.AddVote, api.ValidateAddVote)

	//user route
	v1.Post("/register", userHandler.CreateUser, api.ValidateCreateUser)

	//post route
	v1.Post("/post/:id", postHandler.CreatePost, api.ValidateCreatePost)
	v1.Get("/posts", postHandler.GetPosts)
	v1.Get("/posts/vote", postHandler.GetPostsByVote)
	v1.Get("/posts/order", postHandler.GetPostsByOrder, api.ValidateGetPostsByOrder)

	var postsMutex sync.Mutex
	go func() {
		for {
			time.Sleep(time.Second * 5)
			postsMutex.Lock()
			result, err := config.PostStore.GetPostsByScoreFromPostgres()
			if err != nil {
				log.Fatal(err)
			}
			if (len(result)) != 0 {
				postsJSON, err := json.Marshal(result)
				if err != nil {
					log.Fatal(err)
				}
				err = config.RedisClient.Set(context.Background(), "high_voted_posts", postsJSON, time.Second*5).Err()
			}
			postsMutex.Unlock()
		}
	}()
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", config.ServerPort)))
}
