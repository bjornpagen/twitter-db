package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bjornpagen/twitter-db/db/gen"

	_ "github.com/libsql/libsql-client-go/libsql"
)

//go:generate sqlc generate

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

func (db *DB) AddFullUser(ctx context.Context, u User) error {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	addUserParams := gen.AddUserParams{
		UserID:       u.UserID,
		CreationDate: u.CreationDate,
		Timestamp:    u.Timestamp,
	}
	if err := db.AddUser(ctx, addUserParams); err != nil {
		return fmt.Errorf("add user: %w", err)
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
	if _, err := db.AddUserHistory(ctx, addUserHistoryParams); err != nil {
		return fmt.Errorf("add user history: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
