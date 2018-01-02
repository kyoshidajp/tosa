# TOSA

[![GitHub release](https://img.shields.io/github/release/kyoshidajp/tosa.svg?style=flat-square)][release]
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]

[release]: https://github.com/kyoshidajp/tosa/releases
[license]: https://github.com/kyoshidajp/tosa/blob/master/LICENSE
[godocs]: http://godoc.org/github.com/kyoshidajp/tosa

TOSA is Open Source Available.

Open Pull request page from commit hash. You can easy to find why the code is included by the page:mag_right:

## Usage

```
$ tosa sha
```

*NOTE*: Only first time, `tosa` requires your Github username and password(and two-factor auth code if you are setting). Because of using [GitHub API v3](https://developer.github.com/v3/).

### from tig

Add the following bind command in `~/.tigrc`.

```
bind main O @tosa %(commit)
bind blame O @tosa %(commit)
```

![tig_main_blame](https://user-images.githubusercontent.com/3317191/34467237-ac5e76f4-ef2e-11e7-889d-6d28bf03b04d.gif)

## Install

After install Golang. Then execute `go get`.

```
$ go get -u github.com/kyoshidajp/tosa
```

Or, download binary which meets your system from [Releases](release); then unarchive it and set `$PATH` to the `tosa`.

## Author

Katsuhiko YOSHIDA
