# mynt-backend

## Setting up a development environment
- [Install go](https://golang.org/doc/install)
- [Install sql-migrate](https://github.com/rubenv/sql-migrate)
- [Install golangci-lint](https://pre-commit.com/#install)
- [Install pre-commit](https://pre-commit.com/#install)
- [Install vscode](https://code.visualstudio.com/download)
- [Install the go vscode extension](https://code.visualstudio.com/docs/languages/go)
- [Install docker](https://docs.docker.com/get-docker/)
- Run `pre-commit install` in the root of the repository ([*see this guide for more details*](https://freshman.tech/linting-golang/))
- Create a `.env` file similar to `.env.sample`. `JWT_SIGNING_SECRET` can be whatever you want. For the `GOOGLE_AUTH_CLIENT*` variables, you need to visist the [Google API Console](https://console.developers.google.com/) to obtain Oauth 2.0 credentials.
- Run `docker-compose up` in the root of the repository to run a local development database.
- Run `go run main.go` to start the application.

## Authentication
- We currently use google for authentication.
- To authenticate a client:
    1. Redirect the client to `/redirect` (in this application)
    2. The client will be redirected to the google conscent page
    3. When the client has authenticated on the conscent page, google will redirect the client to `/callback` (in this application)
    4. After this, the `/callback` endpoint will redirect the client to whatever is in the `REDIRECT_TO_FRONTEND` environment variable.
    5. The endpoint also sets a cookie with a jwt-token that can be used to make authenticated requests to this application. The only thing you have to do in javascript is to [set the credentials flag in fetch to include](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch#sending_a_request_with_credentials_included).
- This flow was configured following these guides:
    - https://github.com/coreos/go-oidc
    - https://developers.google.com/identity/protocols/oauth2/openid-connect

## Database migration
- We use [sql-migrate](https://github.com/rubenv/sql-migrate) for managing database updates.
- To make a new database update, run `sql-migrate new replace-with-a-description-of-the-update`.
- The application will make sure that all the migrations are up to date when starting (which means that you don't have to do these updates manually in other environments), but you can run them locally with `sql-migrate up`.
- You can test that the `down` version of your migration works with `sql-migrate redo`.
- Always be carefull with database migrations, we generally don't want to do anything else than creating new tables and new fields. But if you have to, consider doing these things manually instead.

## How to get Sbanken credentials
We currently only support data fetching from the Sbanken Open API platform. To obtain a client id and secret:
- Visit https://secure.sbanken.no/Home/Settings/BetaProgram and enable beta.
- Then, visit https://secure.sbanken.no/Personal/ApiBeta/Info/ to obtain the credentials.
- Note that it often takes some time (~10+ minutes) before theses credentials are valid.

## Required environment variables

- GIN_MODE: release
- PORT: 80
- CALLBACK_URL: https://api.mynt.daae.dev/callback
- REDIRECT_TO_FRONTEND: https://mynt.daae.dev/authenticated/transactions
- REDIRECT_TO_FRONTEND_SIGNOUT: https://mynt.daae.dev
- COOKIE_DOMAIN: daae.dev
- JWT_SIGNING_SECRET
- GOOGLE_AUTH_CLIENT_ID
- GOOGLE_AUTH_CLIENT_SECRET
- POSTGRES_HOST
- POSTGRES_PORT
- POSTGRES_USER
- POSTGRES_PASSWORD
- POSTGRES_DB
- POSTGRES_SSL
