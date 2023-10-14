# gomon

> This started as a fork of
> [JulesGuesnon/Gomon](https://github.com/JulesGuesnon/Gomon)

### Installation

**pre-requisites**

- `curl` if using the [Goblin](#goblin) method
- [go](https://go.dev), assuming this already is installed since you are using a
  go development tool

You can install it using one of the following ways.

#### Goblin

- Request a install script from [goblin.run](http://goblin.run) and follow the
  instructions

```sh
$ curl -sf http://goblin.run/github.com/barelyhuman/gomon@latest | sh
```

#### Go

**Windows**

- Make sure that the `GOPATH` variable is setup and already accessible to be
  able to use `go install`

**Linux/Unix**

- This depends on the `$HOME/go/bin` directory to exist, so make sure that the
  directory exists on Linux/Unix

```
$ go install github.com/barelyhuman/gomon
```

### Examples

A simple project with `nodejs` being used to bundle your frontend assets might
look like so

```sh
$ gomon -i "." -e "./node_modules/*" .
#        ^ include everything from this folder
#               ^ exclude everything from `node_modules` 
#                                    ^ run `go run` on the root of this directory
```

A nested project where the actual binary is somewhere else might look like so

```sh
$ gomon -i "./src" -i "./assets" -e "./node_modules/*" server/bin
#        ^ include everything from this folder to watch
#                   ^ include everything from the assets folder to watch
#                                 ^ exclude everything from the `node_modules` directory
#                                                      ^ run `go run` on the files in `server/bin` directory
```

### CLI

```
NAME:
   gomon - A go program executor with a file watcher

USAGE:
   gomon [global options] command [command options] [arguments...]


COMMANDS:
   watch, w  watch mode
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```
