# STRUCTURE
- All local records are stored in a sqlite file on disk.
- Categories are hierarchical in nature and with a parent/child relationship, can form a graph structure
- The service will be written in go
    - routes should be stored in `src/routes`
    - each route should have its own file, e.g `GET /posts` would be `src/routes/posts`
- the application should have verbose logging throughout execution, both to console and to a file.

# API
- `GET /posts`
- `GET /posts/:category`
- `GET /post/:id`
- `POST /post/new`
- `PUT /post/edit`
- `GET /categories`

# NOTES
`GET /posts` and `GET /posts/:category` should both have parameters of `limit` and `sort` and `offset`. Each response should also include pagination information, of `total_posts`, `total_pages`, `per_page`, and `current_page`. The pages are determined by the `limit` sent through, and the current page is determined via `offset`.
    

# CATEGORIES
These categories should be stored in the sqlite file.

- coding
    - learning
    - devlog
        - gamedev
        - webdev
    - story
- hobbies
    - photography
    - painting
    - hiking
- story
- advice

# TESTING
In order to run the unit tests, the server must first be running, and you must have Bun installed as the test runner.
Simply navigate to the `/tests` folder and run `bun test`. All tests should pass.

If you want logs, run `make 2> <log_file>` and it will redirect stdout to a log file.
