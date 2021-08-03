# mynt-backend

## Data structure

- Accounts
- Transactions
  - Each transaction has a category.
  - Each transaction is linked to an account
- Categories
  - Categories can have parents. So if a transaction has a category, it is also a part of all parents of that category.
- Rules
  - Automatically assigning categories to transactions based on different fields
