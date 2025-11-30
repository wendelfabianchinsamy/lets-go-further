# Migrations

## How are migrations managed?

Migrations are managed using the cli tool golang migration (https://github.com/golang-migrate/migrate).
To install golang-migrate use `brew install golang-migrate`.

## Create a migration

`migrate create -seq -ext=.sql -dir=./migrations create_movies_table`

## What do the different flags indicate in the migration command

- seq: indicates that we want to use sequential numbering like 0001, 0002. Instead of the default unix timestamp
- ext: indicates that we want to give the migration files the extension .sql
- dir: indicates that we want to store the migration files in ./migrations
- create_movies_table: is a descriptive label that we give the migration files to signify their contents

## How to execute migrations>

`migrate -path=./migrations -database='postgres://wendel:password@localhost/letsgofurther?sslmode=disable' up`

We run the migrations "up" on the letsgofurther database.
Notice we single quote the connection string since it has a question mark(?).

## Check which version of migrations your database is running?

You can check the schema_migrations table. It will have a version column indicating the last successful migration run.
You can also check using the migrate CLI tool:
`migrate -path=./migrations -database='postgres://wendel:password@localhost/letsgofurther?sslmode=disable' version`

## Can you migrate up or down to specific version?

Yes you can by running the following migrate CLI command with the "goto" command:

`migrate -path=./migrations -database='postgres://wendel:password@localhost/letsgofurther?sslmode=disable' goto 1`
This will take us to version 1.

`migrate -path=./migrations -database='postgres://wendel:password@localhost/letsgofurther?sslmode=disable' goto 10`
This will take us to version 10.

## Can you migrate down a specific number of migrations?

Yes you can by running the following migrate CLI command with the "down" command:

`migrate -path=./migrations -database='postgres://wendel:password@localhost/letsgofurther?sslmode=disable' down 1`

Here the value after the down command indicates how many versions down you want to migrate to.
It does not represent the version that you want to migrate to.
Running the down command without a value after will rollback all migrations.

## What to do when an error occurs while running a migration?

In the event of an error such as "Dirty database version {X}. Fix and force version." you should take the following steps:

1. Fix the migration code.
2. Force to working version number. If the last correct version was 2 then: `migrate -path=./migrations -database='postgres://wendel:password@localhost/letsgofurther?sslmode=disable' force 2`
3. Run goto for version 2: `migrate -path=./migrations -database='postgres://wendel:password@localhost/letsgofurther?sslmode=disable' goto 2`
4. Now you can re-run your migrations: `migrate -path=./migrations -database='postgres://wendel:password@localhost/letsgofurther?sslmode=disable' up`
