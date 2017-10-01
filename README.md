# Biodiversity Heritage Library Scientific Names Index

Creates an index of scientific names occuring in the collection of literature
in Biodiversity Heritage Library

## Usage

### Mac OSX

* Download [bhlindex release for mac][bhlindex-mac]
* untar the file go to `script` directory and read README.md file

### Linux

* Download [bhlindex release for linux][bhlindex-linux]
* untar the file go to `script` directory and read README.md file

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

[bhlindex-mac]: https://github.com/gnames/bhlindex/releases/download/v0.1.0/bhlindex-0.1.0-mac.tar.gz
[bhlindex-linux]: https://github.com/gnames/bhlindex/releases/download/v0.1.0/bhlindex-0.1.0-linux.tar.gz
