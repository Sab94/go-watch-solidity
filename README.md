# go-watch-solidity

It is a handy little tool to watch a solidity file and generate abi, bin and gobinging.

# Install

`go get -u -x github.com/Sab94/go-watch-solidity`

Note, go-watch-solidity needs `solc` installed to run. Here is the [official guide](https://solidity.readthedocs.io/en/latest/installing-solidity.html) to install solc.

# Usage

```
$ go-watch-solodity -h
   Go Watch Solidity is a watcher for a given solidity file.
   It generates abi, bin, and go bindings for the given solidity
   file on save.

Usage:
  go-watch-solidity [flags]

Flags:
  -a, --abi           Generate abi
  -b, --bin           Generate bin
  -g, --bindgo        Generate go binding (default true)
  -d, --dest string   Destination to generate
  -h, --help          help for go-watch-solidity
```

# Project52

It is one of my [project 52](https://github.com/Sab94/project52).
