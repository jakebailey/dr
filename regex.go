package dr

import "fmt"

// Regex represents the nodes of the regex, which
// support printing a string, as well as transformation
// under derivative.
type Regex interface {
	fmt.Stringer
	Derivative(rune) Regex
	AcceptsEpsilon() bool
}

// Match returns true if the string matches the regex.
// This is sugar for calling Deriative on the regex
// repeatedly, then checking if the current state
// accepts epsilon.
func Match(r Regex, s string) bool {
	for _, c := range s {
		r = r.Derivative(c)
	}
	return r.AcceptsEpsilon()
}

// Empty accepts nothing.
type Empty struct{}

var empty *Empty

func (*Empty) String() string {
	return "∅"
}

// Derivative returns Empty.
func (*Empty) Derivative(rune) Regex {
	return empty
}

// AcceptsEpsilon returns false.
func (*Empty) AcceptsEpsilon() bool {
	return false
}

// Epsilon accepts the empty string.
type Epsilon struct{}

var epsilon *Epsilon

func (*Epsilon) String() string {
	return "ε"
}

// Derivative returns Empty.
func (*Epsilon) Derivative(rune) Regex {
	return empty
}

// AcceptsEpsilon returns true.
func (*Epsilon) AcceptsEpsilon() bool {
	return true
}

// Char only accepts a specific character.
type Char struct {
	R rune
}

var escaped = map[rune]bool{
	'!':  true,
	'(':  true,
	')':  true,
	'*':  true,
	'+':  true,
	'\\': true,
}

func (c *Char) String() string {
	if escaped[c.R] {
		return fmt.Sprintf("\\%c", c.R)
	}
	return fmt.Sprintf("%c", c.R)
}

// Derivative returns Epsilon if r is Char's value,
// otherwise Empty.
func (c *Char) Derivative(r rune) Regex {
	if c.R == r {
		return epsilon
	}
	return empty
}

// AcceptsEpsilon returns false.
func (*Char) AcceptsEpsilon() bool {
	return false
}

// Union accepts the union of two regexes.
type Union struct {
	L Regex
	R Regex
}

func (u *Union) String() string {
	return fmt.Sprintf("(%v)+(%v)", u.L, u.R)
}

// Derivative returns the union of the derivatives
// of this union.
func (u *Union) Derivative(r rune) Regex {
	return &Union{
		L: u.L.Derivative(r),
		R: u.R.Derivative(r),
	}
}

// AcceptsEpsilon returns true if either of the elements
// in the union accepts epsilon.
func (u *Union) AcceptsEpsilon() bool {
	return u.L.AcceptsEpsilon() || u.R.AcceptsEpsilon()
}

// Concat accepts the concatenation of two regexes.
type Concat struct {
	L Regex
	R Regex
}

func (c *Concat) String() string {
	return fmt.Sprintf("%v%v", c.L, c.R)
}

// Derivative returns the union of the concatenation of
// the derivative of L and R, and the derivative of R if
// if L accepts epsilon, otherwise Empty.
func (c *Concat) Derivative(r rune) Regex {
	var right Regex
	if c.L.AcceptsEpsilon() {
		right = c.R.Derivative(r)
	} else {
		right = empty
	}

	return &Union{
		L: &Concat{
			L: c.L.Derivative(r),
			R: c.R,
		},
		R: right,
	}
}

// AcceptsEpsilon returns true if both elements accept epsilon.
func (c *Concat) AcceptsEpsilon() bool {
	return c.L.AcceptsEpsilon() && c.R.AcceptsEpsilon()
}

// Comp accepts the complement of a regex.
type Comp struct {
	R Regex
}

func (c *Comp) String() string {
	return fmt.Sprintf("!(%v)", c.R)
}

// Derivative returns the complement of the derivative of the
// complemented regex.
func (c *Comp) Derivative(r rune) Regex {
	return &Comp{
		R: c.R.Derivative(r),
	}
}

// AcceptsEpsilon returns true if the complemented regex
// does not accept epsilon.
func (c *Comp) AcceptsEpsilon() bool {
	return !c.R.AcceptsEpsilon()
}

// Kleene accepts the Kleene star of a regex.
type Kleene struct {
	R Regex
}

func (k *Kleene) String() string {
	return fmt.Sprintf("(%v)*", k.R)
}

// Derivative returns the concatenation of the derivative
// of the Kleene star'd regex and the regex.
func (k *Kleene) Derivative(r rune) Regex {
	return &Concat{
		L: k.R.Derivative(r),
		R: k,
	}
}

// AcceptsEpsilon returns true.
func (k *Kleene) AcceptsEpsilon() bool {
	return true
}

var (
	_ Regex = &Empty{}
	_ Regex = &Epsilon{}
	_ Regex = &Char{}
	_ Regex = &Union{}
	_ Regex = &Comp{}
	_ Regex = &Kleene{}
)
