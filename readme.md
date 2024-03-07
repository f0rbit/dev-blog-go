# dev-blog
This project is a self-hostable blogging server written in Go meant for small & development-focused blogs. The server itself is written in Go and should be stand-alone. This repo also contains an example client for interacting with the server and serves as examples on how to use the system. The client is a React SPA and as such should be easily deployable anywhere.

## Features
- Mutli-author support using GitHub Authentication
- CRUD operations for blog posts
- Hierarchical categories
- Dynamic tagging system
- Third-party importing from various API
    - Currently supports dev.to

## Installation
### Dependencies
- Makefile
- sqlite3
- go
- node & npm

### Server
After cloning the repo you will want to run `make database` to setup the sqlite database. This will setup the database into `db/sqlite.db` and load all the schemas & migrations. All that's left is making sure the `.env` file is setup correctly and all you have to do is run `make run`. You will need to specify a `PORT` as an environment variable, so from the command line you would run `PORT=8080 make run`.

### Client
To run the client for development you will want to `cd client`, then `npm install` and finally `npm run dev`.
For running the client for deployment you can run `make build-client` from the root directory and then call `node client/dist/entry.mjs` and this will start the production-ready client server.

### .env Setup
For the client the only expected variable in .env file is `VITE_API_URL` which should point to wherever the go server is running. For example if the go server was running on `localhost:8080`, then the `client/.env` file would be
```.env
VITE_API_URL=http://localhost:8080
```
For the server there's a few more variables compared to the client.
```.env
GITHUB_SECRET=<github auth secret>
GITHUB_CLIENT=<github auth client token>
GITHUB_CALLBACK=<go server url>/auth/github/callback
COOKIE_SECRET=<secret hash>
COOKIE_DOMAIN=<go server domain>
CLIENT_URL=<url of client>
```
The `GITHUB_SECRET` and `GITHUB_CLIENT` should be from GitHub's OAuth Integration page which you can find under `Settings` > `Developer Settings` > `OAuth Apps` and after creating a new application, the `GITHUB_CLIENT` will be the `Client ID` and the `GITHUB_SECRET` is under 'Client secrets'.

## Usage
For usage details see the included client as an example.
### Endpoints
| Method | Path                         | Description                                  |
|--------|------------------------------|----------------------------------------------|
| GET    | /posts                       | Fetches all posts.                           |
| GET    | /posts/{category}            | Fetches all posts within a specific category.|
| GET    | /post/{slug}                 | Retrieves a specific post by its slug.       |
| POST   | /post/new                    | Creates a new post.                          |
| PUT    | /post/edit                   | Edits an existing post.                      |
| DELETE | /post/delete/{id}            | Deletes a specific post by its ID.           |
| GET    | /categories                  | Retrieves all categories.                    |
| POST   | /category/new                | Creates a new category.                      |
| DELETE | /category/delete/{name}      | Deletes a specific category by its name.     |
| PUT    | /post/tag                    | Adds a tag to a post.                        |
| DELETE | /post/tag                    | Removes a tag from a post.                   |
| GET    | /tags                        | Retrieves all tags.                          |
| GET    | /auth/user                   | Retrieves information about the logged-in user.|
| GET    | /auth/github/login           | Initiates login via GitHub.                  |
| GET    | /auth/github/callback        | Handles the callback from GitHub authentication.|
| GET    | /auth/logout                 | Logs out the current user.                   |
| GET    | /tokens                      | Retrieves all API tokens for the user.       |
| POST   | /token/new                   | Creates a new API token.                     |
| PUT    | /token/edit                  | Edits an existing API token.                 |
| DELETE | /token/delete/{id}           | Deletes a specific API token by its ID.      |
| GET    | /links                       | Retrieves all integrations for the user.     |
| PUT    | /links/upsert                | Creates or updates an integration.           |
| GET    | /links/fetch/{source}        | Fetches details for a specific integration by source.|
| DELETE | /links/delete/{id}           | Deletes a specific integration by its ID.    |

