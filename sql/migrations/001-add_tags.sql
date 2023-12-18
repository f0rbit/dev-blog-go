-- Create 'tags' table
CREATE TABLE IF NOT EXISTS tags (
    post_id INTEGER,
    tag TEXT,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    UNIQUE (post_id, tag)
);

-- Create index on the 'post_id' column for faster lookups
CREATE INDEX IF NOT EXISTS idx_post_id ON tags(post_id);
