# hashcash

This repo implements [hashcash](http://hashcash.org) in Go. A couple of
simple primitives like `Mint` and `Evaluate` functions are exposed out of
this module.

To use the library, simply import the library from the git repository.

# tests

There are some tests that use the property-based testing
library. However they are not very exhaustive yet.

To run the tests, do `go test . -count=1` (`count=1` disables
successful test results being cached. This is required because we are
using randomized tests that generates different inputs on each test
run)

# License

The code is copyright "Least Authority TFA GmbH" and is licensed under
the MIT License.

