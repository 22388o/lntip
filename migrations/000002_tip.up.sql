CREATE TABLE IF NOT EXISTS tips (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    user_id varchar(255) NOT NULL,
    to_user_id varchar(255) NOT NULL,
    message_id varchar(255),
    guild_id varchar(255) NOT NULL,
    channel_id varchar(255) NOT NULL,
    amount INTEGER NOT NULL,
    is_award BOOLEAN NOT NULL DEFAULT 0,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (to_user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;