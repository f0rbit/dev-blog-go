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
- [x] Expose API to third-party website, make sure it all works
- [x] Add test coverage
    - [x] Tokens
    - [x] Categories
- [x] Deploy client on public-accessible 
    - [x] Try Github pages since there's no need for a server

# v0.9.1
- [x] Description in GET posts,
    - [x] strip markdown from content
    - [x] show first 3 lines

# v1.0
- [x] Integration with third-party API, only support `devto` at the moment, but create a generic interface for adding more
    - [x] devto
- [x] Fetch from these api's cache the results so that when we get a request for /posts, it is near-instant
    - [x] For now, just have a manual 'fetch' button which will ping the API
- [x] Proper README document with instructions on how to self-host/deploy & test locally

# v1.1
- [x] Post Editor redesign
    - [x] Have a dedicated sub-page for editing/creating posts
        - [x] Creation working
        - [x] Editing working
    - [x] Have a 'Preview' section where we can see the content rendered in HTML
        - [x] Add the preview tab & format selection
        - [x] Implement rendering
- [x] Add support for .adoc files
    - [x] Have the filetype configurable per-post
    - [ ] Add support for preview rendering
- [ ] Have a section to write the description for the post in plain text
    - [x] Update description
    - [ ] If left blank, infer description from post

# v1.2
- [ ] Analytics on server
    - [ ] Endpoint for "liking" a post
- [ ] Build homepage
    - [ ] View analytics
    - [ ] Add 'action' table for when someone requests a post (count each request as a 'view')

# v1.3
- [ ] Integrate with media-timeline project
    - [ ] Tightly integrated, not just sharing API keys but full linking workflow that feels natural
    - [ ] Obviously for this, have to finish up media-timeline as well.

# v1.4
- [ ] Light mode support & theme switcher
- [ ] Refactor for support for multiple themes

# Backlog
- [ ] Integrate Medium's API
- [ ] Integrate Substack's API
- [ ] Figure out if we can listen for devto post events and then fetch from there
    - [ ] Otherwise have a task running on the server every hour to refetch all integrations

