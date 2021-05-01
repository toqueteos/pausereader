[![GoDoc][godoc_img]][godoc]
[![Travis Build Status][travis_img]][travis]
[![Go Report Card][goreportcard_img]][goreportcard]
[![codebeat][codebeat_img]][codebeat]

# pausereader

Pausable [io.Reader][ioreader].

**NOTE:** Go modules support was added on `v3`, please check [v3 branch README](https://github.com/toqueteos/pausereader/blob/v3/README.md).

## Installation

```
go get -u github.com/toqueteos/pausereader
```

Project is properly git tagged following [semver][semver] so [glide][glide] and [go dep][go_dep] should be happy.

[godoc_img]: https://godoc.org/github.com/toqueteos/pausereader?status.svg
[godoc]: http://godoc.org/github.com/toqueteos/pausereader
[travis_img]: https://travis-ci.org/toqueteos/pausereader.svg?branch=master
[travis]: https://travis-ci.org/toqueteos/pausereader
[goreportcard_img]: https://goreportcard.com/badge/github.com/toqueteos/pausereader
[goreportcard]: https://goreportcard.com/report/github.com/toqueteos/pausereader
[codebeat_img]: https://codebeat.co/badges/4120eb91-6688-4b9b-9d93-df279a6ebd7f
[codebeat]: https://codebeat.co/projects/github-com-toqueteos-pausereader
[ioreader]: https://golang.org/pkg/io/#Reader
[semver]: http://semver.org/
[glide]: https://glide.sh/
[go_dep]: https://github.com/golang/dep
