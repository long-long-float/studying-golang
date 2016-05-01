(define loop (lambda (n f)
  (cond
    ((eq n nil) nil)
    (t ((lambda () (f) (loop (cdr n) f))))
    )))

(loop (quote 0 0 0 0) (lambda () (print "Hello")))
