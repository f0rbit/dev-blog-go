CREATE TABLE IF NOT EXISTS fetch_queue (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    last_fetch TIMESTAMP NULL,
    location TEXT NOT NULL, -- url
    source TEXT NOT NULL, -- group
    data JSON NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS fetch_links (
    post_id INTEGER NOT NULL,
    fetch_source INTEGER NOT NULL,
    identifier STRING NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (fetch_source) REFERENCES fetch_queue(id),
    UNIQUE(fetch_source, identifier)
);
