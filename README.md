# golisp

(A (lisp) interpreter (in golang)).

Inspired heavily by http://norvig.com/lispy.html. Written on a flight from
Zurich to SFO.

## What it has
* basic math stuff
* lambdas, begin, define, set!, all with proper lexical scoping
* XX? style checks for various bits and pieces
* pretty good error handling (though i started getting lazy with argument count checks)
* test coverage is 60%
* map, car, cdr, etc

## Missing things
* tail-call optimization
* math functions beyond the obvious
* test coverage is only ~60%
* more error handling
* nicer error messages pointng the user to the issues
* cursor navigation in the repl
* brace matching in the repl

## Examples
```lisp
golisp> (define first car)
golisp> (define rest cdr)
golisp> (define count (lambda (item L) (if L (+ (equal? item (first L)) (count item (rest L))) 0)))
golisp> (count 0 (list 0 1 2 3 0 0))
3
golisp> (count (quote the) (quote (the more the merrier the bigger the better)))
4
```
