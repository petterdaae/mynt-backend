# mynt-backend

### Develop
- Create a `.env` file similar to `.env.sample`
- `docker-compose up`
- `go run main.go`

### Authentication
- Authentication is set up with google (create new clients in google console, all you need is a client id and secret)
- Redirect the client to `/redirect`
- You will see the google conscent page
- The backend will redirect to `/authenticated` in the web app and set the `auth_token` cookie that can be used to request protected endpoints.

### Database migration
- `sql-migrate new <name-of-migration>`
- `sql-migrate up`
- `sql-migrate redo`

### Data model

```
User (
  id
  email
)

Account (
  id
  user_id
  external_id
  name
  amount
)

Transaction (
  id
  user_id
  external_id
  account_id
  original_description
  original_date
  amount
  custom_description
  custom_date
  category
)

Category (
  id
  user_id
  name
  parent_category
  monthly_planned_amount
)
```
