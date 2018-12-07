# Biodiversity Heritage Library Indexing Tool binary release

Files in this directory contain binary release of bhlindex tool. The bhlindex
tool finds and records scientific names occuring on >50 million pages
collected by [Biodiversity Heritage Library](http://bhl.org).

## Requirements

1. Laptop or server with >= 8GB of system memory and 500GB of free disk space.
   SSD storage is recommended.

2. Empty postgresql with a user that is able to create new databases.

3. Biodiversity Heritage Library textual files. For testing purposes use files
   located at [testdata directory of this project][testdata]

4. You have to setup environment variables that configure access to BHL files
   and the database server.

    `POSTGRES_DB`
    : Database created for bhlindex

    `POSTGRES_HOST`
    : IP address or hostname where Potgresql database is installed

    `POSTGRES_USER`
    : user that has an access to the POSTGRES_DB

    `BHL_DIR`
    : root of BHL directory that contains `$BHL_DIR`/ocr/bhl1, `$BHL_DIR`/ocr/bhl2 etc.

    `PREF_SOURCES`
    : IDs of data sources from http://resolver.globalnames.org/data_sources.
      They have to be a list of integers separated by comma, for example
      ``PREF_SOURCES=1,2,3``

      The variable with the values for development
      environment can be found at [.env.dev file][env]. To export the variables
      into bash or zsh:

      ```bash
      source .env.dev
      ```

5. Password for the Postgres user should either be empty, or set via
   [`.pgpass` file][pgpass].

## Usage

To check the github commit version and date of compilation use

```bash
./bhlindex version
```

To create the index execute

```bash
./production.sh
```

If you want to read envronment variable from a file

```bash
source /dir/to/env_file ./production.sh
```

[testdata]: https://github.com/gnames/bhlindex/tree/master/testdata/
[env]: https://raw.githubusercontent.com/gnames/bhlindex/master/.env.dev
[pgpass]: https://www.postgresql.org/docs/9.4/static/libpq-pgpass.html
