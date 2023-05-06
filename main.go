package main

import (
	"fmt"
	"log"

	"github.com/JeremyLoy/config"
	twitter "github.com/bjornpagen/rapidapi/twitter154"
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
	tc, err := twitter.New(c.RapidapiKey)
	if err != nil {
		return fmt.Errorf("twitter client: %w", err)
	}

	user, err := tc.GetUserByUsername("bjornpagen")
	if err != nil {
		return fmt.Errorf("get user details: %w", err)
	}

	fmt.Printf("user: %+v\n", user)

	return nil
}
