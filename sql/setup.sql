-- setup.sql

-- Create 'users' table if it doesn't exist
CREATE TABLE IF NOT EXISTS users (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    github_id INTEGER NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    avatar_url VARCHAR(255),
    -- Add other user-related fields as needed
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create 'categories' table if it doesn't exist
CREATE TABLE IF NOT EXISTS categories (
    owner_id INTEGER NOT NULL,
    name VARCHAR(15) PRIMARY KEY,
    parent VARCHAR(15),

    FOREIGN KEY (owner_id) REFERENCES users(user_id),
    UNIQUE(owner_id, name)
);

-- Create 'posts' table if it doesn't exist
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    author_id INTEGER NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT "",
    content TEXT NOT NULL,
    format TEXT NOT NULL DEFAULT "md",
    category TEXT NOT NULL,
    archived BOOLEAN DEFAULT 0,
    publish_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (author_id) REFERENCES users(user_id),
    UNIQUE(author_id, slug)
);

-- create indexes
CREATE INDEX IF NOT EXISTS idx_categories_user ON categories(owner_id);
CREATE INDEX IF NOT EXISTS idx_user_posts ON posts(author_id);

