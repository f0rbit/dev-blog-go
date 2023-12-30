INSERT OR IGNORE INTO users (user_id, github_id, username, email, avatar_url) VALUES
    (1, '81222943', 'f0rbit', 'dev@forbit.dev', 'https://avatars.githubusercontent.com/u/81222943?v=4');

INSERT OR IGNORE INTO categories (owner_id, name, parent) VALUES
    (1, 'coding', 'root'),
    (1, 'learning', 'coding'),
    (1, 'devlog', 'coding'),
    (1, 'gamedev', 'devlog'),
    (1, 'webdev', 'devlog'),
    (1, 'code-story', 'coding'),
    (1, 'hobbies', 'root'),
    (1, 'photography', 'hobbies'),
    (1, 'painting', 'hobbies'),
    (1, 'hiking', 'hobbies'),
    (1, 'story', 'root'),
    (1, 'advice', 'root');

INSERT OR IGNORE INTO posts (slug, author_id, title, content, category) VALUES
    ('test-post', 1, 'test', 'this is a test post, first post.', 'coding');

INSERT OR IGNORE INTO tags (post_id, tag) VALUES
    (1, 'test');

INSERT OR IGNORE INTO access_keys (key_value, user_id) VALUES
    ('hahahahahthisisatestkey', 1);
