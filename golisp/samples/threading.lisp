(define loop (lambda (n f)
  (cond
    ((eq n nil) nil)
    (t ((lambda () (f) (loop (cdr n) f))))
    )))

(define show-10-times (lambda (str)
  (loop (quote 0 0 0 0 0 0 0 0 0 0) (lambda () (print str)))
  ))

(define t1 (thread/run show-10-times "Hello1"))
(define t2 (thread/run show-10-times "Hello2"))

(thread/wait t1)
(thread/wait t2)
