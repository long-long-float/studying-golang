(define get-name (lambda (name)
  (lambda () name)
  ))

(define gn1 (get-name "llf"))
(print (gn1))
