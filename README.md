# golisp

(A (lisp) interpreter (in golang)).

Inspired heavily by http://norvig.com/lispy.html. Written on a flight from
Zurich to SFO.

## What it has
* basic math stuff
* lambdas, begin, define, set!, all with proper lexical scoping
* XX? style checks for various bits and pieces
* pretty good error handling (though i started getting lazy with argument count checks)
* map, car, cdr, etc

## Missing things
* tail-call optimization
* math functions beyond the obvious
* test coverage is ~50%
* more error handling
* nicer error messages pointng the user to the issues
* cursor navigation in the repl
* brace matching in the repl
