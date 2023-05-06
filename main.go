package main

import (
	"fmt"
	"log"
	"time"

	"github.com/JeremyLoy/config"
	twitter "github.com/bjornpagen/rapidapi/twitter154"
	"go.uber.org/ratelimit"
)

type Config struct {
	RapidapiKey string `config:"RAPIDAPI_KEY"`
	DatabaseUrl string `config:"DATABASE_URL"`
}

var c Config

func init() {
	config.FromEnv().To(&c)
}

func validateConfig() {
	unset := make([]string, 0)
	if c.RapidapiKey == "" {
		unset = append(unset, "RAPIDAPI_KEY")
	}
	if c.DatabaseUrl == "" {
		//unset = append(unset, "DATABASE_URL")
	}
	if len(unset) > 0 {
		log.Fatalf("missing required environment variables: %v", unset)
	}
}

func main() {
	validateConfig()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	rl := ratelimit.New(5, ratelimit.Per(time.Second))
	tc, err := twitter.New(c.RapidapiKey, twitter.WithRateLimit(rl))
	if err != nil {
		return fmt.Errorf("twitter client: %w", err)
	}

	user, err := tc.GetUserByUsername("bjornpagen")
	if err != nil {
		return fmt.Errorf("get user details: %w", err)
	}

	userId := user.UserId

	tweets, err := tc.GetUserTweets(userId)
	if err != nil {
		return fmt.Errorf("get user tweets: %w", err)
	}

	for _, tweet := range tweets {
		fmt.Printf("%+v\n", tweet)
	}

	return nil
}
