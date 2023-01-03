# CHANGELOG

## Unreleased

- Add [#61]: shortened and filtered dump for BHL-related data.
- Add [#60]: normalize odds according to verification results.

## [v1.0.0-RC3] - 2022-12-13 Tue

- Add: documentation in README about creating and configuring database and
       its user.
- Fix: show error in case if rebuilding of the database, or its initiation
       did not work as expected.

## [v1.0.0-RC2] - 2022-11-28 Mon

- Add [#59]: modify dump data according to #57.
- Add [#58]: pre-scan data for duplicate page Ids before name-finding.
- Add [#57]: switch to an OCR dump structure that uses item and page Ids.
- Add [#56]: add an option to return all verification results instead
             of the best results only.

## [v1.0.0-RC1] - 2022-10-30 Sun

- Add [#55]: refactor the directory structure using `internal` directory
             to hide code not suitable for public use.

## [v0.13.2] - 2022-09-12 Mon

- Add [#53]: classification ranks and IDs in dump files.

## [v0.13.1] - 2022-09-08 Thu

- Add [#52]: dump pages information

## [v0.13.0] - 2022-09-01 Thu

- Add [#51]: remove RESTful interface, no more remote access.
             All data is taken from dumps.
- Add [#50]: dump saves pages and names separately, allows a flag to
             dump only results for specific data-sources. Dump has a
             flag pointing to a directory where to save dump data.

## [v0.12.6] - 2022-08-29 Mon

- Add: compatibility with GNverifier v1.0.0

## [v0.12.5]

- Add: info for RESTful API.

## [v0.12.4]

- Add: RESTful API for occurrences takes data_sources in account.

## [v0.12.3]

- Add: improve code documentation.
- Add: detected verbatim name to results and data-dump.
- Add: shorten barcode for page to sequence number.
- Fix: deal with verbatim names longer than 255 bytes.

## [v0.12.2]

- Add: improve help messages.

## [v0.12.1]

- Add: update to gnfinder v0.19.2.

## [v0.12.0]

- Add: Update to gnfinder v0.19.1.

## [v0.11.2]

- Fix: add classification ranks, ids to REST API.

## [v0.11.1]

- Add [#49]: add classification ranks, ids.

## [v0.11.0]

- Add [#48]: change RESTful pagination to use IDs.
- Add [#47]: implement `dump` command.
- Add [#45]: create RESTful service.
- Add [#46]: switch to gnverifier for name verification.
- Add [#43]: refactor to improve architecture and usability.

## [v0.10.0]

- Add [#41]: Update to gnfinder v0.11.1.

## [v0.9.0]

- Add [#39]: Save annotations about new species, combinations, subspecies.
- Add [#38]: Save 5 words before and after name-candidates.

## [v0.8.0]

- Add [#36]: Rename `title` to `item` to be in sync with BHL terminology,
  name_string export via gRPC.

## [v0.7.0]

- Add [#35]: Fixes in dictionaries In particular names of botanical genera
  authors are not in the dictionary anymore. Also common latin
  capitalized words from species descriptions are now added to
  'grey' dictionary. As a result calculation of Bayes odds
  score improved quite a bit.
- Add [#34]: There are more indices.
- Add [#32]: Pages are not considered unique anymore and we take a combination
  of item id and archive page id as unique.

## [v0.6.0]

- Add [#31]: save preferred data-sources results to db.
- Add [#30]: average odds and occurrence number for name_strings.
- Add [#29]: matched canonical form from verification.
- Fix [#28]: sporadic non-zero edit distance for ExactMatch.
- Fix [#27]: no verification for abbreviated names.

## [v0.5.0]

- Add [#26]: add Go modules to make builds more stable.
- Add [#24]: updates in verification interface.
- Add [#23]: gRPC has an option to limit stream of pages to specific volumes.
- Add [#22]: gRPC has a stream of volumes metainfo.
- Add [#18]: gRPC example groups names by class clade.
- Add [#17]: gRPC does not stream volumes, streams pages and names and text.
- Add [#16]: gRPC streams volumes, pages, and names.
- Add [#15]: simple gRPC server and an example how to use it.
- Fix [#25]: gRPC serves pages in ascending order instead of random order.

## [v0.4.0]

- Add [#14] curation information for verified names.
- Add [#12],[#13] options to set workers in command line app, better CLI.
- Add [#9],[#10],[#11] improve command line interface.
- Add [#8]: decouple name-finding and name-verification.

## [v0.3.0]

- Add [#4]: set a Makefile to simplify compilation and packaging.
- Add [#3]: verification of name-strings against [gnindex].
- Add [#2]: saving unique name-strings to database.
- Add: gnfinder support for Bayes searches.
- Update: tests to pass again.
- Update: to changes in dependencies.
- Remove: `*.txt` files from `git lfs`.

## [v0.2.0]

- Add: `git lfs` support
- Add: documentation in `README.md` file and script/README.md file.
- Update: to recent `gnfinder`.

## [v0.1.0]

- Add: Biodiversity Heritage Library production trial, 1h for 50 million pages.
- Add: heuristic name finding via gnfinder.
- Add: saving data to database.
- Add: production wrapper script to reset db and do name-finding.
- Add: command line program.
- Add: name-finding framework.
- Add: Postgresql support and migrations.
- Add: development environment with `docker-compose`.

## Footnotes

This document follows [changelog guidelines]

[v1.0.0-RC2]: https://github.com/gnames/bhlindex/compare/v1.0.0-RC1...v1.0.0-RC2
[v1.0.0-RC1]: https://github.com/gnames/bhlindex/compare/v0.13.2...v1.0.0-RC1
[v0.13.2]: https://github.com/gnames/bhlindex/compare/v0.13.1...v0.13.2
[v0.13.1]: https://github.com/gnames/bhlindex/compare/v0.13.0...v0.13.1
[v0.13.0]: https://github.com/gnames/bhlindex/compare/v0.12.6...v0.13.0
[v0.12.6]: https://github.com/gnames/bhlindex/compare/v0.12.5...v0.12.6
[v0.12.5]: https://github.com/gnames/bhlindex/compare/v0.12.4...v0.12.5
[v0.12.4]: https://github.com/gnames/bhlindex/compare/v0.12.3...v0.12.4
[v0.12.3]: https://github.com/gnames/bhlindex/compare/v0.12.2...v0.12.3
[v0.12.2]: https://github.com/gnames/bhlindex/compare/v0.12.1...v0.12.2
[v0.12.1]: https://github.com/gnames/bhlindex/compare/v0.12.0...v0.12.1
[v0.12.0]: https://github.com/gnames/bhlindex/compare/v0.11.2...v0.12.0
[v0.11.2]: https://github.com/gnames/bhlindex/compare/v0.11.1...v0.11.2
[v0.11.1]: https://github.com/gnames/bhlindex/compare/v0.11.0...v0.11.1
[v0.11.0]: https://github.com/gnames/bhlindex/compare/v0.10.0...v0.11.0
[v0.10.0]: https://github.com/gnames/bhlindex/compare/v0.9.0...v0.10.0
[v0.9.0]: https://github.com/gnames/bhlindex/compare/v0.8.0...v0.9.0
[v0.8.0]: https://github.com/gnames/bhlindex/compare/v0.7.0...v0.8.0
[v0.7.0]: https://github.com/gnames/bhlindex/compare/v0.6.0...v0.7.0
[v0.6.0]: https://github.com/gnames/bhlindex/compare/v0.5.0...v0.6.0
[v0.5.0]: https://github.com/gnames/bhlindex/compare/v0.4.0...v0.5.0
[v0.4.0]: https://github.com/gnames/bhlindex/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/gnames/bhlindex/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/gnames/bhlindex/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/gnames/bhlindex/tree/v0.1.0
[#60]: https://github.com/gnames/bhlindex/issues/60
[#59]: https://github.com/gnames/bhlindex/issues/59
[#58]: https://github.com/gnames/bhlindex/issues/58
[#57]: https://github.com/gnames/bhlindex/issues/57
[#56]: https://github.com/gnames/bhlindex/issues/56
[#55]: https://github.com/gnames/bhlindex/issues/55
[#54]: https://github.com/gnames/bhlindex/issues/54
[#53]: https://github.com/gnames/bhlindex/issues/53
[#52]: https://github.com/gnames/bhlindex/issues/52
[#51]: https://github.com/gnames/bhlindex/issues/51
[#50]: https://github.com/gnames/bhlindex/issues/50
[#49]: https://github.com/gnames/bhlindex/issues/49
[#48]: https://github.com/gnames/bhlindex/issues/48
[#47]: https://github.com/gnames/bhlindex/issues/47
[#46]: https://github.com/gnames/bhlindex/issues/46
[#45]: https://github.com/gnames/bhlindex/issues/45
[#44]: https://github.com/gnames/bhlindex/issues/44
[#43]: https://github.com/gnames/bhlindex/issues/43
[#42]: https://github.com/gnames/bhlindex/issues/42
[#41]: https://github.com/gnames/bhlindex/issues/41
[#40]: https://github.com/gnames/bhlindex/issues/40
[#39]: https://github.com/gnames/bhlindex/issues/39
[#38]: https://github.com/gnames/bhlindex/issues/38
[#37]: https://github.com/gnames/bhlindex/issues/37
[#36]: https://github.com/gnames/bhlindex/issues/36
[#35]: https://github.com/gnames/bhlindex/issues/35
[#34]: https://github.com/gnames/bhlindex/issues/34
[#33]: https://github.com/gnames/bhlindex/issues/33
[#32]: https://github.com/gnames/bhlindex/issues/32
[#31]: https://github.com/gnames/bhlindex/issues/31
[#30]: https://github.com/gnames/bhlindex/issues/30
[#29]: https://github.com/gnames/bhlindex/issues/29
[#28]: https://github.com/gnames/bhlindex/issues/28
[#27]: https://github.com/gnames/bhlindex/issues/27
[#26]: https://github.com/gnames/bhlindex/issues/26
[#24]: https://github.com/gnames/bhlindex/issues/24
[#23]: https://github.com/gnames/bhlindex/issues/23
[#22]: https://github.com/gnames/bhlindex/issues/22
[#18]: https://github.com/gnames/bhlindex/issues/18
[#17]: https://github.com/gnames/bhlindex/issues/17
[#16]: https://github.com/gnames/bhlindex/issues/16
[#15]: https://github.com/gnames/bhlindex/issues/15
[#14]: https://github.com/gnames/bhlindex/issues/14
[#13]: https://github.com/gnames/bhlindex/issues/13
[#12]: https://github.com/gnames/bhlindex/issues/12
[#11]: https://github.com/gnames/bhlindex/issues/11
[#10]: https://github.com/gnames/bhlindex/issues/10
[#9]: https://github.com/gnames/bhlindex/issues/9
[#8]: https://github.com/gnames/bhlindex/issues/8
[#4]: https://github.com/gnames/bhlindex/issues/4
[#3]: https://github.com/gnames/bhlindex/issues/3
[#2]: https://github.com/gnames/bhlindex/issues/2
[changelog guidelines]: https://github.com/olivierlacan/keep-a-changelog
[gnindex]: https://index.globalnames.org
