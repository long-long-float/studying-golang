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
* `thread/wait`
  * waits for thread object of argument
