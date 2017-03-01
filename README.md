# dr
ERE derivative-based regex engine

[![GoDoc](https://godoc.org/github.com/jakebailey/dr?status.svg)](http://godoc.org/github.com/jakebailey/dr)

This is based on the derivative-based membership checking from
[CS 598 GR](http://fsl.cs.illinois.edu/index.php/CS598_-_Runtime_Verification_(Spring_2017)), aka Runtime Verification.

The implementation is relatively trivial (it's quite literally the same as the paper),
so it might as well be public. The parser was the time consuming part.

"dr" comes from dR, i.e. derivative of R, because I think I'm funny.

The syntax uses `!` for complement, `*` for the Kleene star, `+` for union, and
`.` for any character.

The output of the test program in `cmd/drtest` is:

```
asdfg => asdfg
aaa+bbb => (aaa)+(bbb)
!(a)b*(cd)*e+f => (!(a)(b)*(cd)*e)+(f)
\+\++\*(\!\\) => (\+\+)+(\*\!\\)

matching against ab(c)*
: false
a: false
ab: true
abc: true
abccccc: true
```

Which shows the input and output after parsing and generating the regex, as well
as various examples of matching a common expression.

In addition to those rules described in class, I've also added a rule for `.` (any),
which accepts any single character, although, it could have been represented by
`(!a+a)`.