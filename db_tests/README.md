# Testing over DB (PostgreSQL) with Goose experiments

There are three moments about testing over DB presented here.

## Stair test.

Represented in `repository_migrate`.
Strategy:
1. Parse migrations with
```go
migrations, err = goose.CollectMigrations(migrationsDir, 0, goose.MaxVersion)
```
2. Create empty DB.
3. Generate tests for each migration with up-down-up circle.

If no migration fail - we can assume, that there are backwards migrations for all migrations. 
All migrations may be applied forwards and backwards without database corruption.
Stair test is not enough to check if there is no data corruption - because one need specific tests for such cases.

## Repository test over DB with DB copies.

Represented in `repository/names`.

Strategy:
1. Create empty DB.
2. Apply all migrations.
3. Make this DB a template.
4. Create new DB from template for each test. Use it and drop it.
5. Make template DB a common DB again and drop it.

This strategy guaranties absolute isolation between tests. It also allows to test complex things that include 
transactions, rollbacks and so on. You have right the same behavior in database as in production.
Drawback of this strategy is that it is slow.

## Repository test over DB with transactions.

Represented in `repository/properties`.

Strategy:
1. Create empty DB.
2. Apply all migrations.
3. Create transaction for each test. Rollback it after test.
4. Drop DB.

This strategy may cause problems with transactions. For example, if your repository works with rollbacks somehow for 
example. Or you can catch problems like represented in `repository/properties` test `Should update property updated_at` 
that forced me to use `statement_timestamp()` instead of `now()`. But this strategy is fast in my tests it showed to 
be two orders faster than previous one. In cases with massive migrations the difference will be even greater.

## Existing problems

Both cases may be improved if there is some coordinator that guaranties DB with migrations. One cannot make this 
optimization without coordinator, because tests in different packages are independent and may not know about each other.
