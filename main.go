package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/JeremyLoy/config"
	twitter "github.com/bjornpagen/rapidapi/twitter154"
	"go.uber.org/ratelimit"

	"github.com/bjornpagen/twitter-db/db"
	"github.com/bjornpagen/twitter-db/db/gen"
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
		unset = append(unset, "DATABASE_URL")
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

	myDB, err := db.New(c.DatabaseUrl)
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}

	user, err := tc.GetUserByUsername("redacted")
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	// add user to db
	if _, err = myDB.AddFullUser(context.Background(), toDBUser(user)); err != nil {
		return fmt.Errorf("add user to db: %w", err)
	}

	err = fetchTweets(&tc, &myDB, user.UserId)
	if err != nil {
		return fmt.Errorf("fetch tweets: %w", err)
	}

	ids, err := myDB.GetTweetIDs(context.Background(), user.UserId)
	if err != nil {
		return fmt.Errorf("get tweets: %w", err)
	}

	for _, id := range ids {
		// fetch likers of Tweet
		users, err := tc.GetTweetUserFavorites(id)
		if err != nil {
			return fmt.Errorf("get tweet user favorites: %w", err)
		}

		// add likers to db
		ctx := context.Background()
		for _, user := range users {
			if _, err = myDB.AddFullUser(ctx, toDBUser(user)); err != nil {
				return fmt.Errorf("add user to db: %w", err)
			}

			if err = myDB.AddFavorite(ctx, gen.AddFavoriteParams{
				UserID:  user.UserId,
				TweetID: id,
			}); err != nil {
				return fmt.Errorf("add user favorite: %w", err)
			}
		}
	}

	return nil
}

func fetchTweets(tc *twitter.Client, myDB *db.DB, userId string) error {
	tweets, err := tc.GetUserTweets(userId, twitter.IncludePinned())
	if err != nil {
		return fmt.Errorf("get user tweets: %w", err)
	}

	ctx := context.Background()
	for _, tweet := range tweets {
		if _, err = myDB.AddFullTweet(ctx, toDBTweet(tweet, userId)); err != nil {
			return fmt.Errorf("add user to db: %w", err)
		}
	}

	return nil
}

type number interface {
	float32 | float64 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

func to[T number](n bool) T {
	if n {
		return 1
	}
	return 0
}

func toDBUser(u twitter.User) db.User {
	return db.User{
		UserID:           u.UserId,
		CreationDate:     u.CreationDate,
		Timestamp:        int64(u.Timestamp),
		Username:         u.Username,
		Name:             u.Name,
		FollowerCount:    int64(u.FollowerCount),
		FollowingCount:   int64(u.FollowingCount),
		FavouritesCount:  int64(u.FavouritesCount),
		IsPrivate:        to[int64](u.IsPrivate),
		IsVerified:       to[int64](u.IsVerified),
		IsBlueVerified:   to[int64](u.IsBlueVerified),
		Location:         u.Location,
		ProfilePicUrl:    u.ProfilePicUrl,
		ProfileBannerUrl: u.ProfileBannerUrl,
		Description:      u.Description,
		ExternalUrl:      u.ExternalUrl,
		NumberOfTweets:   int64(u.NumberOfTweets),
		Bot:              to[int64](u.Bot),
		HasNftAvatar:     to[int64](u.HasNftAvatar),
		DefaultProfile:   to[int64](u.DefaultProfile),
		DefaultImage:     to[int64](u.DefaultImage),
	}
}

func toDBTweet(t twitter.Tweet, userId string) db.Tweet {
	return db.Tweet{
		TweetID:        t.TweetId,
		UserID:         userId,
		CreationDate:   t.CreationDate,
		Text:           t.Text,
		MediaUrl:       t.MediaUrl,
		VideoUrl:       toDBVideoUrl(t.VideoUrl),
		Language:       t.Language,
		FavoriteCount:  int64(t.FavoriteCount),
		RetweetCount:   int64(t.RetweetCount),
		ReplyCount:     int64(t.ReplyCount),
		QuoteCount:     int64(t.QuoteCount),
		Retweet:        to[int64](t.Retweet),
		Views:          int64(t.Views),
		Timestamp:      int64(t.Timestamp),
		VideoViewCount: int64(t.VideoViewCount),
		ExpandedUrl:    t.ExpandedUrl,
		ConversationID: t.ConversationId,
	}
}

func toDBVideoUrl(vu []twitter.VideoUrl) []db.VideoUrl {
	v := make([]db.VideoUrl, len(vu))
	for i, url := range vu {
		v[i] = db.VideoUrl{
			Bitrate:     int64(url.Bitrate),
			ContentType: url.ContentType,
			Url:         url.Url,
		}
	}
	return v
}

/*
users, err := tc.GetTweetUserFavorites("1639998740462223363")
	if err != nil {
		return fmt.Errorf("get tweet user favorites: %w", err)
	}

	spew.Dump(users)

	os.Exit(0)

	// add likers to db
	ctx := context.Background()
	for _, user := range users {
		if _, err = myDB.AddFullUser(ctx, toDBUser(user)); err != nil {
			return fmt.Errorf("add user to db: %w", err)
		}

		if err = myDB.AddFavorite(ctx, gen.AddFavoriteParams{
			UserID:  user.UserId,
			TweetID: "1639998740462223363",
		}); err != nil {
			return fmt.Errorf("add user favorite: %w", err)
		}
	}
*/
