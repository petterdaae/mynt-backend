# mynt-backend

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
- `/sync`
  - `POST` pull new transactions from banks
