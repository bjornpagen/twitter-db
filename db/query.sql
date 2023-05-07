-- name: GetUser :one
SELECT * FROM users
WHERE user_id = ? LIMIT 1;

-- name: AddUser :exec
INSERT INTO users (user_id)
VALUES (?);

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = ?;

-- name: GetUserHistory :one
SELECT * FROM user_history
WHERE user_id = ? ORDER BY timestamp DESC LIMIT 1;

-- name: AddUserHistory :one
INSERT INTO user_history (
	creation_date,
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
	timestamp,
	has_nft_avatar,
	default_profile,
	default_image
) VALUES (
	?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING id;

-- name: GetFollowers :many
SELECT follower_id FROM follow
WHERE user_id = ?;

-- name: GetFollowing :many
SELECT user_id FROM follow
WHERE follower_id = ?;

-- name: AddFollow :exec
INSERT INTO follow (
	user_id,
	follower_id,
	timestamp
) VALUES (
	?, ?, ?
);

-- name: DeleteFollow :exec
DELETE FROM follow
WHERE user_id = ? AND follower_id = ?;

-- name: GetTweet :one
SELECT * FROM tweets
WHERE tweet_id = ? LIMIT 1;

-- name: AddTweet :exec
INSERT INTO tweets (tweet_id, user_id)
VALUES (?, ?);

-- name: DeleteTweet :exec
DELETE FROM tweets
WHERE tweet_id = ?;

-- name: GetTweetHistory :one
SELECT * FROM tweet_history
WHERE tweet_id = ? ORDER BY timestamp DESC LIMIT 1;

-- name: AddTweetHistory :one
INSERT INTO tweet_history (
	creation_date,
	tweet_id,
	text,
	user_id,
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
) VALUES (
	?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING id;