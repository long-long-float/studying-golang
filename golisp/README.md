# Lisp.go

An pure lisp intepreter implemented by Go.

## Build

```sh
$ go build
```

## Usage

```sh
$ ./golisp FILE
```

# Spec

* functions
  * `atom`, `eq`, `car`, `cdr`, `cons`
  * `print`
* special forms
  * `cond`, `quote`, `lambda`, `define`
* values
  * integer(`"-"? [1-9][0-9]*`), char(`"'" [^'] "'"`), string(list of char `"\"" [^"]* "\""`)
