#BHL Index

Creates an Index of scientific names mentioned in the collection of
Biodiversity Heritage Library

## Database Migrations

### Install migrate

```bash
go get -u -d github.com/mattes/migrate/cli github.com/lib/pq
go build -tags 'postgres' -o $GOPATH/bin/migrate github.com/mattes/migrate/cli
```

### Create migration

```
migrate -ext sql -D db NAME
```

### Run commands

```
migrate -database postgres://localhost:5432/database up 2
```

### Commands

create [-ext E] [-dir D] NAME
: Create a set of timestamped up/down migrations titled NAME, in
  directory D with extension E

version
: current migration version

up [N]
: up N migrations

down [N]
: down N migrations

drop
: nuke database

