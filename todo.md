# v0.1
- [x] `GET /posts`
- [x] `GET /posts/:cateogry`
- [x] `GET /post/:id`
- [x] `POST /post/new`
- [x] `PUT /post/edit`
- [x] `GET /categories`
- [x] `DELETE /post/delete/:id`
- [x] Primary testing of endpoints

# v0.2
- [ ] `GET /tags`
- [ ] `PUT /post/tag`
- [ ] `GET /posts` (with tag param)
- [x] Move all SQL setup & seeding queries into .sql files
- [x] setup a migrations folder (for tags)
- [x] Replace logging with Charm's Go logging library

# v0.3
- [ ] Deployment scripts & setup
- [ ] Github actions for testing
- [ ] Figure out how to get code coverage over the two completely different codebases
- [ ] Add authentication middleware for endpoints

# v0.4
- [ ] Front end implementation
    - [ ] Create Blog Posts
    - [ ] Edit existing Posts
    - [ ] Delete posts
    - [ ] Assign categories
    - [ ] Assign tags
- [ ] Front end typesafety w/ zod
- [ ] Front end styling

# v0.5
- [ ] Database connection pooling
- [ ] Proper graph data structure for categories
- [ ] Refactor & simplify code bases
    - [ ] > 85% test coverage
- [ ] Proper input/output data validation (type safety)
