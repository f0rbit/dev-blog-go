INSERT OR IGNORE INTO users (user_id, github_id, username, email, avatar_url) VALUES
    (1, '81222943', 'f0rbit', 'dev@forbit.dev', 'https://avatars.githubusercontent.com/u/81222943?v=4');

-- Insert data into 'posts' table if it's empty
INSERT OR IGNORE INTO posts (slug, author_id, title, content, category) VALUES
    ('test-post', 1, 'test', 'this is a test post, first post.', 'coding');

INSERT OR IGNORE INTO tags (post_id, tag) VALUES
    (1, 'test');
