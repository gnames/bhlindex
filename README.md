# Biodiversity Heritage Library Scientific Names Index

Creates an index of scientific names occurring in the collection of literature
in Biodiversity Heritage Library

## Installation on Linux

- Download [bhlindex release for Linux][bhlindex-linux]
- Untar the file, copy it to `/usr/local/bin` or other directory in the `PATH`.
- Use [bhl testdata][bhl-test] for testing.

## Configuration

When you run the app for the first time it will create a configuration file and
will provide information where the file is located (usually it is
`$HOME/.config/bhlnames.yaml`)

Edit the file to provide credentials for PostgreSQL database.

Change the `Jobs` setting according to the amount of memory and the number
of CPU. For 32Gb of memory `Jobs: 8` works ok. This parameter sets the number
of concurrent jobs running for name-finding.

Set `BHLdir` parameter to point to the root directory where BHL texts are
located (several hundred gigabytes of texts).

Other parameters a optional.

## Usage

### Preparations

Login to PostgreSQL server and create a database that has the same name as the
`PgDatabase` parameter in the configuration file (default name is `bhlindex`).

This database will be used to keep found names. Its final size of the database
upon completion should be in a vicinity of 50GB.

### Commands

Find names in BHL

```bash
bhlindex find
```

Verify detected names using [GNverifier] service

```bash
bhlindex verify
```

Dump data into tab-separated files

```bash
bhlindex dump
```

Serve detected items, pages, verified names, names occurrences vie RESTful
interface (default port is 8080).

```bash
bhlindex rest
# using different port
bhlindex rest -p 8000
```

To run name detection, verification and dump as one command:

```bash
bhlindex all
```

### Testing

```bash
docker-compose build
docker-compose up
```

To update all dependencies change LAST_FULL_REBUILD line in Docker file and
return `docker-compose build`

[bhlindex-mac]: https://github.com/gnames/bhlindex/releases/download/v0.1.0/bhlindex-0.1.0-mac.tar.gz
[bhlindex-linux]: https://github.com/gnames/bhlindex/releases/download/v0.1.0/bhlindex-0.1.0-linux.tar.gz
[bhl-test]: https://github.com/gnames/bhlindex/releases/download/v0.1.0/bhl-testdata.tar.gz
[readme]: https://github.com/gnames/bhlindex/tree/master/bhlindex
