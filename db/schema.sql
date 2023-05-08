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
	is_private INTEGER NOT NULL CHECK(is_private IN (0, 1)),
	is_verified INTEGER NOT NULL CHECK(is_verified IN (0, 1)),
	is_blue_verified INTEGER NOT NULL CHECK(is_blue_verified IN (0, 1)),
	location TEXT NOT NULL,
	profile_pic_url TEXT NOT NULL,
	profile_banner_url TEXT NOT NULL,
	description TEXT NOT NULL,
	external_url TEXT NOT NULL,
	number_of_tweets INTEGER NOT NULL,
	bot INTEGER NOT NULL CHECK(bot IN (0, 1)),
	has_nft_avatar INTEGER NOT NULL CHECK(has_nft_avatar IN (0, 1)),
	default_profile INTEGER NOT NULL CHECK(default_profile IN (0, 1)),
	default_image INTEGER NOT NULL CHECK(default_image IN (0, 1)),
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

CREATE TABLE favorite (
	user_id TEXT NOT NULL,
	tweet_id TEXT NOT NULL,
	row_created INTEGER NOT NULL DEFAULT (unixepoch('now')),

	PRIMARY KEY (user_id, tweet_id),
	FOREIGN key(user_id) REFERENCES users(user_id),
	FOREIGN key(tweet_id) REFERENCES tweets(tweet_id)
) STRICT;

CREATE TABLE retweet (
	user_id TEXT NOT NULL,
	tweet_id TEXT NOT NULL,
	row_created INTEGER NOT NULL DEFAULT (unixepoch('now')),

	PRIMARY KEY (user_id, tweet_id),
	FOREIGN key(user_id) REFERENCES users(user_id),
	FOREIGN key(tweet_id) REFERENCES tweets(tweet_id)
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
	text TEXT NOT NULL,
	language TEXT NOT NULL,
	favorite_count INTEGER NOT NULL,
	retweet_count INTEGER NOT NULL,
	reply_count INTEGER NOT NULL,
	quote_count INTEGER NOT NULL,
	retweet INTEGER NOT NULL CHECK(retweet IN (0, 1)),
	views INTEGER NOT NULL,
	timestamp INTEGER NOT NULL,
	video_view_count INTEGER NOT NULL,
	expanded_url TEXT NOT NULL,
	conversation_id TEXT NOT NULL,
	FOREIGN KEY (tweet_id) REFERENCES tweets(tweet_id),
	FOREIGN KEY (user_id) REFERENCES users(user_id)
) STRICT;

CREATE INDEX tweet_history_tweet_id_idx ON tweet_history(tweet_id);
CREATE INDEX tweet_history_user_id_idx ON tweet_history(user_id);
CREATE INDEX tweet_history_timestamp_idx ON tweet_history(timestamp);

CREATE TABLE media_urls (
	tweet_history_id INTEGER NOT NULL,
	url TEXT NOT NULL,
	row_created INTEGER NOT NULL DEFAULT (unixepoch('now')),

	PRIMARY KEY (tweet_history_id, url),
	FOREIGN KEY (tweet_history_id) REFERENCES tweet_history(id)
) STRICT;

CREATE TABLE video_urls (
	tweet_history_id INTEGER NOT NULL,
	url TEXT NOT NULL,
	row_created INTEGER NOT NULL DEFAULT (unixepoch('now')),

	bitrate INTEGER NOT NULL,
	content_type TEXT NOT NULL,
	PRIMARY KEY (tweet_history_id, url),
	FOREIGN KEY (tweet_history_id) REFERENCES tweet_history(id)
) STRICT;