# Running
To run the server, simply run `make run`. You need to support the `PORT` as an environment variable, and for authentication middleware you can provide `AUTH_TOKEN` environment variable as well. If you want logs, run `make run 2> <log_file>` and it will redirect stdout to a log file.

# Structure
- All local records are stored in a sqlite file on disk.
- Categories are hierarchical in nature and with a parent/child relationship, can form a graph structure
- The service will be written in go
    - routes should be stored in `src/routes`
    - each route should have its own file, e.g `GET /posts` would be `src/routes/posts`
- the application should have verbose logging throughout execution, both to console and to a file.

# API
- `GET /posts`
- `GET /posts/:category`
- `GET /post/:slug`
- `POST /post/new`
- `PUT /post/edit`
- `DELETE /post/delete/:id`
- `PUT /post/tag`
- `DELETE /post/tag`
- `GET /categories`
- `GET /tags`

# Parameters
`GET /posts` and `GET /posts/:category` should both have parameters of `limit` and `sort` and `offset`. Each response should also include pagination information, of `total_posts`, `total_pages`, `per_page`, and `current_page`. The pages are determined by the `limit` sent through, and the current page is determined via `offset`.
    

# Categories
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

# Testing
To run the integration tests run `make test`. If you want to see line coverage run `make coverage`. The tests are written in TypeScript and uses Bun as the runtime for fast & efficient testing.

