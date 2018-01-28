# TOSA

[![GitHub release](https://img.shields.io/github/release/kyoshidajp/tosa.svg?style=flat-square)][release]
[![Travis](https://travis-ci.org/kyoshidajp/tosa.svg?branch=master)](https://travis-ci.org/kyoshidajp/tosa)
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]

[release]: https://github.com/kyoshidajp/tosa/releases
[license]: https://github.com/kyoshidajp/tosa/blob/master/LICENSE
[godocs]: http://godoc.org/github.com/kyoshidajp/tosa

**TOSA** is **O**pen **S**ource **A**vailable.

Open pull request page or get pull request data from sha(commit hash). You can more easily find why the code is included by the page:mag_right:

![tosa](https://user-images.githubusercontent.com/3317191/35214690-57fabb5e-ffa5-11e7-90b0-1a1d9f50ebab.gif)

## Usage

```
$ tosa <sha>
```

### Options

```
-u, --url      Print the pull request url.

-a, --apiurl   Print the issue API url.

-n, --newline  If -u(--url) or -a(--apiurl) option is specified, print
               the url with newline character at last.

-d, --debug    Enable debug mode.
               Print debug log.

-h, --help     Show this help message and exit.

-v, --version  Print current version.
```

*NOTE*: Set Github Access Token which has "Full control of private repositories" scope as an environment variable `GITHUB_TOKEN`. If not set, `tosa` requires your Github username and password(and two-factor auth code if you are setting). Because of using [GitHub API v3](https://developer.github.com/v3/).


### from tig

Add the following key bindings in `$HOME/.tigrc`.

```
bind main O @tosa %(commit)
bind blame O @tosa %(commit)
```

Open page by O(Shift+o) in main or blame view.

![tig_main_blame](https://user-images.githubusercontent.com/3317191/34467237-ac5e76f4-ef2e-11e7-889d-6d28bf03b04d.gif)

### API URL

Get GitHub issue API url.

```
$ tosa -a <sha>
```

Get title of pull request via [jq](https://stedolan.github.io/jq/), for example.

```
$ curl -s `tosa -a c97e6909` | jq -r '.title'
Add short command option and usage
```

For more information, see [Issues \| GitHub Developer Guide](https://developer.github.com/v3/issues/#get-a-single-issue). 

## Install

### Homebrew

If you have already installed [Homebrew](http://brew.sh/); then can install by brew command.

```
$ brew tap kyoshidajp/tosa
$ brew install tosa
```

### go get

If you are a Golang developper/user; then execute `go get`.

```
$ go get -u github.com/kyoshidajp/tosa
```

### Manual

1. Download binary which meets your system from [Releases](release).
1. Unarchive it.
1. Put `tosa` where you want.
1. Add `tosa` path to `$PATH`.

## Configuration

### Browser

Add `browser` in `$HOME/.config/tosa`(which is YAML format). Like this.

```
github.com:
- user: your_account
  oauth_token: your_oauth_token
  protocol: https
  browser: firefox
```

This config file is automatically generated when running without setting environment variable `GITHUB_TOKEN`.

The value is a name or an absolute path. By default, used your system default browser.

## Author

[Katsuhiko YOSHIDA](https://github.com/kyoshidajp)
