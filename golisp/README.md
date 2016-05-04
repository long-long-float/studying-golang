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

## Spec

* functions
  * `atom`, `eq`, `car`, `cdr`, `cons`
  * `print`
  * `+`, `-`, `*`, `/`, `%` (available for only integer)
* special forms
  * `cond`, `quote`, `lambda`, `define`
* values
  * integer(`"-"? [1-9][0-9]*`), char(`"'" [^'] "'"`), string(list of char `"\"" [^"]* "\""`)

## Threading

```lisp
(define loop (lambda (n f)
  (cond
    ((eq n nil) nil)
    (t ((lambda () (f) (loop (cdr n) f))))
    )))

(define show-10-times (lambda ()
  (loop (quote 0 0 0 0 0 0 0 0 0 0) (lambda () (print "Hello")))
  ))

(define t1 (thread/run show-10-times))
(define t2 (thread/run show-10-times))

(thread/wait t1)
(thread/wait t2)
```

* `thread/run`
  * starts running thread, and returns thread object
  * Program doesn't stopping even if an error happens at thread, so you can use `thread/wait` to receive an error.
* `thread/wait`
  * waits for thread object, and returns result of thread or an error.
