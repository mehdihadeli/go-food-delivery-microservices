# Migration Using Go-Migrate
- [migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [Getting started](https://github.com/golang-migrate/migrate/blob/856ea12df9d230b0145e23d951b7dbd6b86621cb/GETTING_STARTED.md)
- [Migration Filename Format](https://github.com/golang-migrate/migrate/blob/856ea12df9d230b0145e23d951b7dbd6b86621cb/MIGRATIONS.md)
- [PostgreSQL tutorial for beginners](https://github.com/golang-migrate/migrate/blob/856ea12df9d230b0145e23d951b7dbd6b86621cb/database/postgres/TUTORIAL.md)

# Atlas Go-Migrate
[https://atlasgo.io/guides/migration-tools/golang-migrate](https://atlasgo.io/guides/migration-tools/golang-migrate)

```bash
atlas migrate hash --env go-migrate
```

```bash
atlas schema inspect --env go-migrate --url "file://db/migrations/go-migrate" --format "{{ sql . \"  \" }}" > ./db/migrations/go-migrate/schema.sql
```
