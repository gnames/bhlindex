# Biodiversity Heritage Library Scientific Names Index

Creates an index of scientific names occurring in the collection of literature
in Biodiversity Heritage Library

## Performance

This application allows to traverse all digitized corpus of Biodiversity
Heritage Library in a matter of hours. On a modern high-end laptop we
observed the following results:

- name-finding in 275,000 volumes, 60 million pages: `2.5 hours`.
- name-verification of 23 million unique name-strings: `3 hours`.
- preparing a CSV file with 250 million names occurrences/verification records
  : `40 minutes`.

## Installation on Linux

- Download [bhlindex latest release for Linux][bhlindex-latest]
- Untar the file, copy it to `/usr/local/bin` or other directory in the `PATH`.
- Use [bhl testdata][bhl-test] for testing.

BHL corpus of OCRed data can be found as a [>50GB compressed file][bhl-ocr].

## Configuration

When you run the app for the first time it will create a configuration file and
will provide information where the file is located (usually it is
`$HOME/.config/bhlnames.yaml`)

Edit the file to provide credentials for PostgreSQL database.

Change the `Jobs` setting according to the amount of memory and the number
of CPU. For 32Gb of memory `Jobs: 7` works ok. This parameter sets the number
of concurrent jobs running for name-finding.

Set `BHLdir` parameter to point to the root directory where BHL texts are
located (several hundred gigabytes of texts).

Other parameters are optional.

### Environment Variables

It is possible to use Environment Variables instead of configuration file.
Environment Variable override the configuration file settings. The following
variable can be used:

| Config         | Env. Variable        |
| -------------- | -------------------- |
| BHLdir         | BHLI_BHL_DIR         |
| OutputFormat   | BHLI_OUTPUT_FORMAT   |
| PgHost         | BHLI_PG_HOST         |
| PgPort         | BHLI_PG_PORT         |
| PgUser         | BHLI_PG_USER         |
| PgPass         | BHLI_PG_PASS         |
| PgDatabase     | BHLI_PG_DATABASE     |
| Jobs           | BHLI_JOBS            |
| VerifierURL    | BHLI_VERIFIER_URL    |
| WithWebLogs    | BHLI_WITH_WEB_LOGS   |
| WithoutConfirm | BHLI_WITHOUT_CONFIRM |

## Usage

### Preparations

Login to PostgreSQL server and create a database that has the same name as the
`PgDatabase` parameter in the configuration file (default name is `bhlindex`).

This database will be used to keep found names. Its final size of the database
upon completion should be in a vicinity of 50GB.

### Commands

Get BHLindex version

```bash
bhlindex -V
```

Find names in BHL

```bash
bhlindex find
# to avoid confirmation dialog (-y overrides configuration file)
bhlindec find -y
```

Verify detected names using [GNverifier] service

```bash
bhlindex verify
# to avoid confirmation dialog (-y overrides configuration file)
bhlindec verify -y
```

Dump data into tab-separated files

```bash
bhlindex dump
# to compress and save on disk
bhlindex dump | gzip > bhlindex-dump.csv.gz

# -f overrides configuration file settings for output format
bhlindex dump -f tsv | gzip > bhlindex-dump.tsv.gz
bhlindex dump -f json | gzip > bhlindex-dump.json.gz
```

To run all commands together

```bash
bhlindex find -y && \
  bhlindex verify -y && \
  bhlindex dump | gzip > bhlindex-dump.csv.gz
```

Serve detected items, pages, verified names, names occurrences via RESTful
interface (default port is 8080).

```bash
bhlindex rest
# using different port
bhlindex rest -p 8000
```

## RESTful API endpoints

- `/api/v0/items`
- `/api/v0/pages`
- `/api/v0/names`
- `/api/v0/occurrences`

| Query                                         | Usage                                                                 |
| --------------------------------------------- | --------------------------------------------------------------------- |
| items?offset_id=11&limit=100                  | get items with ids 11-110                                             |
| pages?offset_id=11&limit=10                   | get pages of items with ids 11-20                                     |
| names?offset_id=1&limit=10                    | get verified names with ids 1-10                                      |
| names?offset_id=1&limit=10&data_sources=1     | get verified names with ids 1-10 verified to the "Catalogue of Life"  |
| occurrences?offset=21&limit=10                | get detected names with ids 21-30                                     |
| occurrences?offset=21&limit=10&data_sources=1 | get detected names with ids 21-30 verified to the "Catalogue of Life" |

### Testing

Testing requires PostgreSQL database `bhlindex_test`.
Testing will delete all data from the database.

```bash
go test
```

[bhl-ocr]: http://opendata.globalnames.org/dumps/
[bhlindex-latest]: https://github.com/gnames/bhlindex/releases/latest
[bhl-test]: https://github.com/gnames/bhlindex/tree/master/testdata/bhl/ocr
[readme]: https://github.com/gnames/bhlindex/tree/master/bhlindex
