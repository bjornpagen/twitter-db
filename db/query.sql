-- name: GetUser :one
SELECT * FROM users
WHERE user_id = ? LIMIT 1;

-- name: AddUser :exec
INSERT OR IGNORE INTO users (user_id, creation_date, timestamp)
VALUES (?, ?, ?);

-- name: GetLatestUserHistory :one
SELECT * FROM user_history
WHERE user_id = ?
ORDER BY row_created DESC LIMIT 1;

-- name: AddUserHistory :one
INSERT INTO user_history (
	user_id,
	username,
	name,
	follower_count,
	following_count,
	favourites_count,
	is_private,
	is_verified,
	is_blue_verified,
	location,
	profile_pic_url,
	profile_banner_url,
	description,
	external_url,
	number_of_tweets,
	bot,
	has_nft_avatar,
	default_profile,
	default_image
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id;

-- name: AddFollow :exec
INSERT OR REPLACE INTO follow (user_id, follower_id)
VALUES (?, ?);

-- name: AddTweet :exec
INSERT OR IGNORE INTO tweets (tweet_id, user_id)
VALUES (?, ?);

-- name: GetLatestTweetHistory :one
SELECT * FROM tweet_history
WHERE tweet_id = ?
ORDER BY timestamp DESC LIMIT 1;

-- name: AddTweetHistory :one
INSERT INTO tweet_history (
	tweet_id,
	user_id,
	creation_date,
	text,
	language,
	favorite_count,
	retweet_count,
	reply_count,
	quote_count,
	retweet,
	views,
	timestamp,
	video_view_count,
	expanded_url,
	conversation_id
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id;

-- name: GetFollowers :many
SELECT * FROM follow
WHERE user_id = ?;

-- name: GetFollowing :many
SELECT * FROM follow
WHERE follower_id = ?;

-- name: GetTweets :many
SELECT * FROM tweets
WHERE user_id = ?;
