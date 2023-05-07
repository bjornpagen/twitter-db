package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bjornpagen/twitter-db/db/gen"

	_ "github.com/libsql/libsql-client-go/libsql"
)

//go:generate sqlc generate

var (
	ErrNotImplemented = fmt.Errorf("not implemented")
)

type DB struct {
	gen.Queries
	db *sql.DB
}

type User struct {
	UserID           string
	CreationDate     string
	Timestamp        int64
	Username         string
	Name             string
	FollowerCount    int64
	FollowingCount   int64
	FavouritesCount  int64
	IsPrivate        int64
	IsVerified       int64
	IsBlueVerified   int64
	Location         string
	ProfilePicUrl    string
	ProfileBannerUrl string
	Description      string
	ExternalUrl      string
	NumberOfTweets   int64
	Bot              int64
	HasNftAvatar     int64
	DefaultProfile   int64
	DefaultImage     int64
}

type Tweet struct {
	TweetID        string
	UserID         string
	CreationDate   string
	Text           string
	MediaUrl       []string
	VideoUrl       []VideoUrl
	Language       string
	FavoriteCount  int64
	RetweetCount   int64
	ReplyCount     int64
	QuoteCount     int64
	Retweet        int64
	Views          int64
	Timestamp      int64
	VideoViewCount int64
	ExpandedUrl    string
	ConversationID string
}

type VideoUrl struct {
	Bitrate     int64
	ContentType string
	Url         string
}

func New(dbUrl string) (DB, error) {
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		return DB{}, fmt.Errorf("open database connection: %w", err)
	}

	return DB{
		*gen.New(db),
		db,
	}, nil
}

func (db *DB) AddFullUser(ctx context.Context, u User) (userHistoryId int64, err error) {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	addUserParams := gen.AddUserParams{
		UserID:       u.UserID,
		CreationDate: u.CreationDate,
		Timestamp:    u.Timestamp,
	}
	if err := db.AddUser(ctx, addUserParams); err != nil {
		return 0, fmt.Errorf("add user: %w", err)
	}

	addUserHistoryParams := gen.AddUserHistoryParams{
		UserID:           u.UserID,
		Username:         u.Username,
		Name:             u.Name,
		FollowerCount:    u.FollowerCount,
		FollowingCount:   u.FollowingCount,
		FavouritesCount:  u.FavouritesCount,
		IsPrivate:        u.IsPrivate,
		IsVerified:       u.IsVerified,
		IsBlueVerified:   u.IsBlueVerified,
		Location:         u.Location,
		ProfilePicUrl:    u.ProfilePicUrl,
		ProfileBannerUrl: u.ProfileBannerUrl,
		Description:      u.Description,
		ExternalUrl:      u.ExternalUrl,
		NumberOfTweets:   u.NumberOfTweets,
		Bot:              u.Bot,
		HasNftAvatar:     u.HasNftAvatar,
		DefaultProfile:   u.DefaultProfile,
		DefaultImage:     u.DefaultImage,
	}
	if userHistoryId, err = db.AddUserHistory(ctx, addUserHistoryParams); err != nil {
		return 0, fmt.Errorf("add user history: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit transaction: %w", err)
	}

	return userHistoryId, nil
}

func (db *DB) AddFullTweet(ctx context.Context, t Tweet) (tweetHistoryId int64, err error) {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	addTweetParams := gen.AddTweetParams{
		TweetID: t.TweetID,
		UserID:  t.UserID,
	}
	if err := db.AddTweet(ctx, addTweetParams); err != nil {
		return 0, fmt.Errorf("add tweet: %w", err)
	}

	addTweetHistoryParams := gen.AddTweetHistoryParams{
		TweetID:        t.TweetID,
		UserID:         t.UserID,
		CreationDate:   t.CreationDate,
		Text:           t.Text,
		Language:       t.Language,
		FavoriteCount:  t.FavoriteCount,
		RetweetCount:   t.RetweetCount,
		ReplyCount:     t.ReplyCount,
		QuoteCount:     t.QuoteCount,
		Retweet:        t.Retweet,
		Views:          t.Views,
		Timestamp:      t.Timestamp,
		VideoViewCount: t.VideoViewCount,
		ExpandedUrl:    t.ExpandedUrl,
		ConversationID: t.ConversationID,
	}
	if tweetHistoryId, err = db.AddTweetHistory(ctx, addTweetHistoryParams); err != nil {
		return 0, fmt.Errorf("add tweet history: %w", err)
	}

	for _, v := range t.VideoUrl {
		addVideoUrlParams := gen.AddVideoUrlParams{
			TweetHistoryID: tweetHistoryId,
			Bitrate:        v.Bitrate,
			ContentType:    v.ContentType,
			Url:            v.Url,
		}
		if err := db.AddVideoUrl(ctx, addVideoUrlParams); err != nil {
			return 0, fmt.Errorf("add video url: %w", err)
		}
	}

	for _, m := range t.MediaUrl {
		addMediaUrlParams := gen.AddMediaUrlParams{
			TweetHistoryID: tweetHistoryId,
			Url:            m,
		}
		if err := db.AddMediaUrl(ctx, addMediaUrlParams); err != nil {
			return 0, fmt.Errorf("add media url: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit transaction: %w", err)
	}

	return tweetHistoryId, nil
}

func (db *DB) GC(ctx context.Context) error {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := db.GCUsers(ctx); err != nil {
		return fmt.Errorf("gc users: %w", err)
	}

	if err := db.GCTweets(ctx); err != nil {
		return fmt.Errorf("gc tweets: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (db *DB) GCUsers(ctx context.Context) error {
	return ErrNotImplemented
}

func (db *DB) GCTweets(ctx context.Context) error {
	return ErrNotImplemented
}
