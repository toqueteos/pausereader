[![GoDoc][godoc_img]][godoc]
![GHA Build Status][gha_img]
[![Go Report Card][goreportcard_img]][goreportcard]
[![codebeat][codebeat_img]][codebeat]

# pausereader

Pausable [io.Reader][ioreader].

## Installation

```bash
go get github.com/toqueteos/pausereader/v3
```

## Changelog

- `v3` Go modules support
- `v2` updated library internals to `sync` & `sync/atomic` primitives instead of timer-based polling
- `v1` first release (timer-based polling)

[godoc_img]: https://godoc.org/github.com/toqueteos/pausereader?status.svg
[godoc]: http://godoc.org/github.com/toqueteos/pausereader
[gha_img]: https://github.com/toqueteos/pausereader/actions/workflows/test.yml/badge.svg
[goreportcard_img]: https://goreportcard.com/badge/github.com/toqueteos/pausereader
[goreportcard]: https://goreportcard.com/report/github.com/toqueteos/pausereader
[codebeat_img]: https://codebeat.co/badges/4120eb91-6688-4b9b-9d93-df279a6ebd7f
[codebeat]: https://codebeat.co/projects/github-com-toqueteos-pausereader
[ioreader]: https://golang.org/pkg/io/#Reader
