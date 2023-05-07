CREATE TABLE users (
	user_id TEXT PRIMARY KEY,
	row_created INTEGER NOT NULL DEFAULT (unixepoch('now')),

	creation_date TEXT NOT NULL,
	timestamp INTEGER NOT NULL
) STRICT;

CREATE TABLE user_history (
	id INTEGER PRIMARY KEY,
	user_id TEXT NOT NULL,
	row_created INTEGER NOT NULL DEFAULT (unixepoch('now')),

	username TEXT NOT NULL,
	name TEXT NOT NULL,
	follower_count INTEGER NOT NULL,
	following_count INTEGER NOT NULL,
	favourites_count INTEGER NOT NULL,
	is_private INTEGER NOT NULL,
	is_verified INTEGER NOT NULL,
	is_blue_verified INTEGER NOT NULL,
	location TEXT NOT NULL,
	profile_pic_url TEXT NOT NULL,
	profile_banner_url TEXT NOT NULL,
	description TEXT NOT NULL,
	external_url TEXT NOT NULL,
	number_of_tweets INTEGER NOT NULL,
	bot INTEGER NOT NULL,
	has_nft_avatar INTEGER NOT NULL,
	default_profile INTEGER NOT NULL,
	default_image INTEGER NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(user_id)
) STRICT;

CREATE INDEX user_history_user_id_idx ON user_history(user_id);
CREATE INDEX user_history_row_created_idx ON user_history(row_created);

CREATE TABLE follow (
	user_id TEXT NOT NULL,
	follower_id TEXT NOT NULL,
	row_created INTEGER NOT NULL DEFAULT (unixepoch('now')),

	PRIMARY KEY (user_id, follower_id),
	FOREIGN key(user_id) REFERENCES users(user_id),
	FOREIGN key(follower_id) REFERENCES users(user_id)
) STRICT;

CREATE TABLE tweets (
	tweet_id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL,
	row_created INTEGER NOT NULL DEFAULT (unixepoch('now')),

	FOREIGN key (user_id) REFERENCES users(user_id)
) STRICT;

CREATE INDEX tweets_user_id_idx ON tweets(user_id);

CREATE TABLE tweet_history (
	id INTEGER PRIMARY KEY,
	tweet_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	row_created INTEGER NOT NULL DEFAULT (unixepoch('now')),

	creation_date TEXT NOT NULL,
	text TEXT,
	language TEXT,
	favorite_count INTEGER,
	retweet_count INTEGER,
	reply_count INTEGER,
	quote_count INTEGER,
	retweet INTEGER,
	views INTEGER,
	timestamp INTEGER,
	video_view_count INTEGER,
	expanded_url TEXT,
	conversation_id TEXT,
	FOREIGN KEY (tweet_id) REFERENCES tweets(tweet_id),
	FOREIGN KEY (user_id) REFERENCES users(user_id)
) STRICT;

CREATE INDEX tweet_history_tweet_id_idx ON tweet_history(tweet_id);
CREATE INDEX tweet_history_user_id_idx ON tweet_history(user_id);
CREATE INDEX tweet_history_timestamp_idx ON tweet_history(timestamp);

CREATE TABLE media_urls (
	id INTEGER PRIMARY KEY,
	tweet_id TEXT NOT NULL,
	row_created INTEGER NOT NULL DEFAULT (unixepoch('now')),

	url TEXT NOT NULL,
	FOREIGN KEY (tweet_id) REFERENCES tweet_history(tweet_id)
) STRICT;

CREATE INDEX media_urls_tweet_id_idx ON media_urls(tweet_id);

CREATE TABLE video_urls (
	id INTEGER PRIMARY KEY,
	tweet_id TEXT NOT NULL,
	row_created INTEGER NOT NULL DEFAULT (unixepoch('now')),

	bitrate INTEGER NOT NULL,
	content_type TEXT NOT NULL,
	url TEXT NOT NULL,
	FOREIGN KEY (tweet_id) REFERENCES tweet_history(tweet_id)
) STRICT;

CREATE INDEX video_urls_tweet_id_idx ON video_urls(tweet_id);