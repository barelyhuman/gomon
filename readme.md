# Gomon

> This is fork of [JulesGuesnon/Gomon](https://github.com/JulesGuesnon/Gomon)

This package aim to reproduce the behavior of [nodemon](https://github.com/remy/nodemon) for go.
I made this for training purpose so it's probably not really usable.

## Installation guide

Install the package

```sh
go install github.com/barelyhuman/gomon
```

You can also run this directly using

```sh
go get -u github.com/barelyhuman/gomon
go run github.com/barelyhuman/gomon <flags and options>
```

There you go !

## How to use it ?

For now you can only watch a file, nothing else

```sh
gomon path/to/my/file.go
# or
gomon -w "./src,./dist" path/to/my/file.go
```

## Possible issue

If you face this issue:

```
gomon: command not found
```

You may need to add `GOPATH` to your `PATH` (you may need to set your `GOPATH`)

```sh
export PATH=$PATH:$GOPATH/bin
```
