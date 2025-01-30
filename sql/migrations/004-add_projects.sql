-- tables for proejcts, we also want to cache projects that we fetch from devpad_api
CREATE TABLE IF NOT EXISTS projects_cache (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    fetched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status TEXT NOT NULL, -- pending, fetching, fetched
    url TEXT NOT NULL,
    data JSON NULL
);

-- because we can't append project_id to posts table, we will create a table for connecting projects to posts, storing the project_uuid and project_id
CREATE TABLE IF NOT EXISTS posts_projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    project_uuid TEXT NOT NULL,
    project_id INTEGER NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(id)
);
