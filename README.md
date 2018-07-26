# Biodiversity Heritage Library Scientific Names Index

Creates an index of scientific names occuring in the collection of literature
in Biodiversity Heritage Library

## Usage

*NOTE*: This is a Dev release.

### Mac OSX

* Download [bhlindex release for mac][bhlindex-mac].
* Untar the file go to `script` directory and [read instructions][readme].
* Use [bhl testdata][bhl-test] for testing.

### Linux

* Download [bhlindex release for linux][bhlindex-linux]
* Untar the file go to `script` directory and [read instructions][readme].
* Use [bhl testdata][bhl-test] for testing.

## Database Migrations

### Install migrate

```bash
go get -u -d github.com/golang-migrate/migrate/cli github.com/lib/pq
go build -tags 'postgres' -o $GOPATH/bin/migrate github.com/golang-migrate/migrate/cli
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

### Testing

```
docker-compose build
docker-compose up
```

To update all dependencies change LAST_FULL_REBUILD line in Docker file and
return `docker-compose build`


[bhlindex-mac]: https://github.com/gnames/bhlindex/releases/download/v0.1.0/bhlindex-0.1.0-mac.tar.gz
[bhlindex-linux]: https://github.com/gnames/bhlindex/releases/download/v0.1.0/bhlindex-0.1.0-linux.tar.gz
[bhl-test]: https://github.com/gnames/bhlindex/releases/download/v0.1.0/bhl-testdata.tar.gz
[readme]: https://github.com/gnames/bhlindex/tree/master/bhlindex
