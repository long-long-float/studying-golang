(define isNil (lambda (x) (eq x nil)))

(print (isNil 1))
(print (isNil nil))
(print (isNil (atom (quote))))
