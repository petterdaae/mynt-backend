# mynt-backend

- Authentication is set up with auth0. I followed [this guies](https://auth0.com/docs/quickstart/backend/golang#validate-access-tokens) when setting things up.

### Data model

```
Account (
  id
  external_id
  name
  amount
)

Transaction (
  id
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
  name
  parent_category
  monthly_planned_amount
)
```

### Endpoints

- `/accounts`
  - `GET` all accounts
  - `POST` insert new account
- `/accounts/<account_id>`
  - `PUT` update account
- `/transactions?startDate=2020-01-01&endDate=2021-01-01`
  - `GET` all transactions
  - `POST` insert new transaction
- `/transactions/<transaction_id>`
  - `PUT` update transaction, should only be possible to change custom_description, custom_date and category
- `/categories`
  - `GET` all categories
  - `POST` insert new category
- `/categories/<category_id>`
  - `PUT` update category
