# Biodiversity Heritage Library Indexing Tool binary release.

Files in this directory contain binary a release of bhlindex tool to
run reindexing of scientific names that occure on >50 million pages collected
by Biodiversity Heritage Library (http://bhl.org)

## Requirements

1. Laptop or server with >= 8GB of sytem memory and 500GB of free disk space.
   SSD storage is recommended.

2. Empty postgresql with a user that is able to creaate new databases.

3. Biodiversity Heritage Library textual files. For testing purposes use files
   located at [testdata directory of this project][testdata]

4. You have to setup environment variables that configure access to BHL files
   and the database server. The variable with the values for development
   environment can be found at [.env.dev file][env]. Password for the Postgres
   user should either be empty, or setup using [`.pgpass` file][pgpass].

## Usage

To check the github commit version and date of compilation use

```
./bhlindex version
```

To create the index execute

```
./production.sh
```
If you want to read envronment variable from a file

```
source /dir/to/env_file ./production.sh
```

[testdata]: https://github.com/gnames/bhlindex/tree/master/testdata/
[env]: https://raw.githubusercontent.com/gnames/bhlindex/master/.env.dev
[pgpass]: https://www.postgresql.org/docs/9.4/static/libpq-pgpass.html

