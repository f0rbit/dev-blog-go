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
- [x] `GET /tags`
- [x] `PUT /post/tag`
- [x] `DELETE /post/tag`
- [x] `GET /posts` (with tag param)
- [x] Move all SQL setup & seeding queries into .sql files
- [x] setup a migrations folder (for tags)
- [x] Replace logging with Charm's Go logging library

# v0.3
- [x] Deployment scripts & setup
- [x] Github actions for testing
- [x] Figure out how to get code coverage over the two completely different codebases
- [x] Add authentication middleware for endpoints
- [x] Update README for server

# v0.4
- [ ] Front end implementation
    - [x] Create Blog Posts
    - [x] Edit existing Posts
    - [x] Delete posts
    - [x] Assign categories
    - [x] Assign tags
    - [x] Filter posts menu
        - [x] By category
        - [x] By tag
- [ ] Front end typesafety w/ zod
- [x] Front end styling
- [ ] Implement publish date on posts

# v0.5
- [ ] Database connection pooling
- [ ] Proper graph data structure for categories
- [ ] Refactor & simplify code bases
    - [ ] > 85% test coverage
- [ ] Proper input/output data validation (type safety)

# v0.6
- [ ] Instant Reactivity in client, consider transitioning to UUID for id's in DB
    - [ ] Blog updates & creation should appear instant in client, some kind of spinner and validation message upon completion
- [ ] Tag Input component like category component
- [ ] Arrow key down in dropdowns should scroll the div down 
- [ ] Login Page
    - [ ] Proper authentication token handling

# v1.0
- [ ] Expose API to third-party website
- [ ] Deploy client on public-accessible 
    - [ ] Try Github pages since there's no need for a server
    - [ ] Otherwise try vercel (free tier)
    - [ ] Finally, host client on vps alongside the server

# v1.1
- [ ] Transition DB into a multi-user setup
- [ ] Allow for multiple 'authors' on the server
- [ ] Authentication support for each author individually
- [ ] Client supports multiple 'authors'

# v1.2
- [ ] Integration with third-party API, only support `devto` at the moment, but create a generic interface for adding more
    - [ ] devto
    - [ ] Medium's API
    - [ ] Substack's API
- [ ] Integrate with media-timeline project
    - [ ] Tightly integrated, not just sharing API keys but full linking workflow that feels natural
    - [ ] Obviously for this, have to finish up media-timeline as well.

# v1.3
- [ ] Light mode support & theme switcher
- [ ] Refactor for support for multiple themes

