-- Insert data into 'posts' table if it's empty
INSERT OR IGNORE INTO posts (slug, title, content, category) VALUES
    ('test-post', 'test', 'this is a test post, first post.', 'coding');
