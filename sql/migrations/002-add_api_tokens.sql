CREATE TABLE IF NOT EXISTS access_keys (
    key_id INTEGER PRIMARY KEY AUTOINCREMENT,
    key_value TEXT NOT NULL,
    user_id INTEGER NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(user_id),
    UNIQUE(key_value)
);

CREATE INDEX IF NOT EXISTS idx_key_value ON access_keys(key_value);
