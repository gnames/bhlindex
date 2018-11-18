# Changelog

## Unreleased

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

[v0.5.0]: https://github.com/gnames/bhlindex/compare/v0.4.0...v0.5.0
[v0.4.0]: https://github.com/gnames/bhlindex/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/gnames/bhlindex/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/gnames/bhlindex/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/gnames/bhlindex/tree/v0.1.0

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
