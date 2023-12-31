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
- [x] Front end implementation
    - [x] Create Blog Posts
    - [x] Edit existing Posts
    - [x] Delete posts
    - [x] Assign categories
    - [x] Assign tags
    - [x] Filter posts menu
        - [x] By category
        - [x] By tag
- [x] Front end styling
- [x] Implement publish date on posts

# v0.5
- [x] Database connection pooling
- [x] Proper graph data structure for categories
- [x] Refactor & simplify code bases
    - [x] Back end
    - [x] Front end
- [x] > 80% test coverage
- [x] Front end typesafety w/ zod
    - [x] Create zod schemas
    - [x] Validate responses

# v0.6
- [x] Instant Reactivity in client
    - [x] Blog updates & creation should appear instant in client, spinner while request is pending
- [x] Tag Input component like category component
- [x] Arrow key down in dropdowns should scroll the div down 
- [x] Login Page
    - [x] Password input
    - [x] Route for testing token
    - [x] Store token and include in headers to server

# v0.7
- [x] Transition DB into a multi-user setup
- [x] Allow for multiple 'authors' on the server
- [x] Client supports multiple 'authors'
- [x] Authentication support for each author individually
    - [x] Backend logging in
    - [x] Front end log in page / logout button
- [x] Fix testing with new authentication flow
    - [x] Create a "test" user
    - [x] Have an api-keys table that allows passing through a token in header to access as a user

# v0.8
- [x] Client implementation of managing api tokens in settings
    - [x] Create
    - [x] Edit
    - [x] Delete
- [x] Server access tokens CRUD operations
    - [x] Create token
    - [x] Update existing token
    - [x] Delete tokens
    - [x] Get tokens
- [x] Create interface for adding/removing categories for each user
    - [x] Create category
    - [x] Remove category
    - [x] Get categories

# v0.9
- [ ] Expose API to third-party website, make sure it all works
- [ ] Add test coverage
    - [x] Tokens
    - [ ] Categories
- [x] Deploy client on public-accessible 
    - [x] Try Github pages since there's no need for a server

# v1.0
- [ ] Integration with third-party API, only support `devto` at the moment, but create a generic interface for adding more
    - [ ] devto
    - [ ] Medium's API
    - [ ] Substack's API
- [ ] Integrate with media-timeline project
    - [ ] Tightly integrated, not just sharing API keys but full linking workflow that feels natural
    - [ ] Obviously for this, have to finish up media-timeline as well.

# v1.1
- [ ] Analytics on server
    - [ ] Endpoint for "liking" a post
- [ ] Build homepage
    - [ ] View analytics

# v1.2
- [ ] Light mode support & theme switcher
- [ ] Refactor for support for multiple themes

