# Biodiversity Heritage Library Scientific Names Index (BHLindex)

[![Doc Status][doc-img]][doc]

Creates an index of scientific names occurred in the collection of literature
in the Biodiversity Heritage Library

## Performance

This application allows to traverse all digitized corpus of Biodiversity
Heritage Library in a matter of hours. On a modern high-end laptop we
observed the following results:

- name-finding in 275,000 volumes, 60 million pages: `2.5 hours`.
- name-verification of 23 million unique name-strings: `3 hours`.
- preparing a CSV file with 250 million names occurrences/verification records
  : `40 minutes`.

## Installation on Linux

- Download [the bhlindex's latest release for Linux][bhlindex-latest]
- Untar the file, copy it to `/usr/local/bin` or other directory in the `PATH`.
- Use [bhl testdata][bhl-test] for testing.

BHL corpus of OCR-ed data can be found as a [>50GB compressed file][bhl-ocr].

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
Environment Variables override the configuration file settings. The following
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

Three files will be created: `pages`, `names`, `occurrences`. They
will have extension according to selected output format (CSV is the default).
If it is required to filter verified results by data-sources, their list and
corresponding IDs can be found at [gnverifier sources page]

Uncompressed dump files take more than 30GB of space.

```bash
# Dump files to a designated directory.
bhlindex dump -d ~/bhlindex-dump
# or
bhlindex dump --dir ~/bhlindex-dump

# Dump records verified to particular data-sources of `gnverifier`.
# In this case verified names are filtered by `The Catalogue of Life` (ID=1)
# and `The Encyclopedia of Life` (ID=12).
bhlindex dump -d ~/bhlindex-dump -s 1,12
or
bhlindex dump --dir ~/bhlindex-dump --sources 1,12

# Dump using JSON or TSV formats.
bhlindex dump -f tsv -d ~/bhlindex-dump
bhlindex dump -f json -d ~/bhlindex-dump
#or
bhlindex dump --format tsv --dir ~/bhlindex-dump
```

To run all commands together

```bash
bhlindex find -y && \
  bhlindex verify -y && \
  bhlindex dump -d output-dir
```

### Testing

Testing requires PostgreSQL database `bhlindex_test`.
Testing will delete all data from the test database.

```bash
go test ./...
```

[bhl-ocr]: http://opendata.globalnames.org/dumps/
[bhlindex-latest]: https://github.com/gnames/bhlindex/releases/latest
[bhl-test]: https://github.com/gnames/bhlindex/tree/master/testdata/bhl/ocr
[doc-img]: https://godoc.org/github.com/gnames/bhlindex?status.png
[doc]: https://godoc.org/github.com/gnames/bhlindex
