-- setup.sql

-- Create 'categories' table if it doesn't exist
CREATE TABLE IF NOT EXISTS categories (
    name VARCHAR(15) PRIMARY KEY,
    parent VARCHAR(15)
);

-- Insert data into 'categories' table if it's empty
INSERT OR IGNORE INTO categories (name, parent) VALUES
    ('coding', 'root'),
    ('learning', 'coding'),
    ('devlog', 'coding'),
    ('gamedev', 'devlog'),
    ('webdev', 'devlog'),
    ('code-story', 'coding'),
    ('hobbies', 'root'),
    ('photography', 'hobbies'),
    ('painting', 'hobbies'),
    ('hiking', 'hobbies'),
    ('story', 'root'),
    ('advice', 'root');

-- Create 'posts' table if it doesn't exist
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    category TEXT NOT NULL,
    archived BOOLEAN DEFAULT 0,
    publish_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

