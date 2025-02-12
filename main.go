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
	)

	//vote route
	app.Post("/vote/:id", voteHandler.AddVote, api.ValidateAddVote)

	//user route
	app.Post("/register", userHandler.CreateUser, api.ValidateCreateUser)

	//post route
	app.Post("/post/:id", postHandler.CreatePost, api.ValidateCreatePost)
	app.Get("/posts", postHandler.GetPosts)
	app.Get("/posts/vote", postHandler.GetPostsByVote)
	app.Get("/posts/order", postHandler.GetPostsByOrder, api.ValidateGetPostsByOrder)

	var postsMutex sync.Mutex
	go func() {
		for {
			time.Sleep(time.Second * 5)
			postsMutex.Lock()
			result, _ := config.PostStore.GetPostsByScoreFromPostgres()
			if (len(result)) != 0 {
				postsJSON, _ := json.Marshal(result)
				err = config.RedisClient.Set(context.Background(), "high_voted_posts", postsJSON, time.Second*5).Err()
			}
			postsMutex.Unlock()
		}
	}()
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", config.ServerPort)))
}
