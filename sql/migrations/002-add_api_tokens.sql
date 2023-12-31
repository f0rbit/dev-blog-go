CREATE TABLE IF NOT EXISTS access_keys (
    key_id INTEGER PRIMARY KEY AUTOINCREMENT,
    key_value TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    note TEXT NOT NULL,
    enabled BOOLEAN DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(user_id),
    UNIQUE(key_value)
);

CREATE INDEX IF NOT EXISTS idx_key_value ON access_keys(key_value);
