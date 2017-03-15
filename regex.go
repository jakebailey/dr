package dr

import "fmt"

// Regex represents the nodes of the regex, which
// support printing a string, as well as transformation
// under derivative.
type Regex interface {
	fmt.Stringer
	Derivative(rune) Regex
	Accepting() bool
}

// Match returns true if the string matches the regex.
// This is sugar for calling Derivative on the regex
// repeatedly, then checking if the current state
// accepts epsilon.
func Match(r Regex, s string) bool {
	for _, c := range s {
		r = r.Derivative(c)
	}
	return r.Accepting()
}

type empty struct{}

// NewEmpty creates a regex that accepts nothing.
func NewEmpty() Regex {
	return (*empty)(nil)
}

func (*empty) String() string {
	return "∅"
}

// Derivative returns empty.
func (*empty) Derivative(rune) Regex {
	return NewEmpty()
}

// Accepting returns false.
func (*empty) Accepting() bool {
	return false
}

type epsilon struct{}

// NewEpsilon creates a regex that only accepts the empty string.
func NewEpsilon() Regex {
	return (*epsilon)(nil)
}

func (*epsilon) String() string {
	return "ε"
}

// Derivative returns empty.
func (*epsilon) Derivative(rune) Regex {
	return NewEmpty()
}

// Accepting returns true.
func (*epsilon) Accepting() bool {
	return true
}

type char struct {
	r rune
}

// NewChar creates a regex that only accepts the given character.
func NewChar(r rune) Regex {
	return &char{r: r}
}

var escaped = map[rune]bool{
	'!':  true,
	'(':  true,
	')':  true,
	'*':  true,
	'+':  true,
	'\\': true,
	'.':  true,
}

func (c *char) String() string {
	if escaped[c.r] {
		return fmt.Sprintf("\\%c", c.r)
	}
	return fmt.Sprintf("%c", c.r)
}

// Derivative returns Epsilon if r is Char's value,
// otherwise Empty.
func (c *char) Derivative(r rune) Regex {
	if c.r == r {
		return NewEpsilon()
	}
	return NewEmpty()
}

// Accepting returns false.
func (*char) Accepting() bool {
	return false
}

type any struct{}

// NewAny creates a regex that accepts any single character.
func NewAny() Regex {
	return (*any)(nil)
}

func (*any) String() string {
	return "."
}

// Derivative returns epsilon.
func (*any) Derivative(rune) Regex {
	return NewEpsilon()
}

// Accepting return false.
func (*any) Accepting() bool {
	return false
}

// Union accepts the union of two regexes.
type union struct {
	l Regex
	r Regex
}

// NewUnion creates a regex that accepts the union of two regexes,
// taking into consideration the simplifiying equations.
func NewUnion(l, r Regex) Regex {
	switch l.(type) {
	case *empty:
		return r
	default:
		return &union{
			l: l,
			r: r,
		}
	}
}

func (u *union) String() string {
	return fmt.Sprintf("(%v)+(%v)", u.l, u.r)
}

// Derivative returns the union of the derivatives
// of this union.
func (u *union) Derivative(r rune) Regex {
	return NewUnion(u.l.Derivative(r), u.r.Derivative(r))
}

// Accepting returns true if either of the elements
// in the union are accepting.
func (u *union) Accepting() bool {
	return u.l.Accepting() || u.r.Accepting()
}

type concat struct {
	l Regex
	r Regex
}

// NewConcat creates a regex that accepts the concatenation of two regexes,
// taking into consideration the simplifiying equations.
func NewConcat(l, r Regex) Regex {
	switch l.(type) {
	case *empty:
		return NewEmpty()
	case *epsilon:
		return r
	default:
		return &concat{
			l: l,
			r: r,
		}
	}
}

func (c *concat) String() string {
	return fmt.Sprintf("%v%v", c.l, c.r)
}

// Derivative returns the union of the concatenation of
// the derivative of L and R, and the derivative of R if
// if L accepts epsilon, otherwise Empty.
func (c *concat) Derivative(r rune) Regex {
	var right Regex
	if c.l.Accepting() {
		right = c.r.Derivative(r)
	} else {
		right = NewEmpty()
	}

	return NewUnion(
		NewConcat(c.l.Derivative(r), c.r),
		right,
	)
}

// Accepting returns true if both elements are accepting.
func (c *concat) Accepting() bool {
	return c.l.Accepting() && c.r.Accepting()
}

type comp struct {
	r Regex
}

// NewComp creates a regex that accepts the complement of a regex.
func NewComp(r Regex) Regex {
	return &comp{
		r: r,
	}
}

func (c *comp) String() string {
	return fmt.Sprintf("!(%v)", c.r)
}

// Derivative returns the complement of the derivative of the
// complemented regex.
func (c *comp) Derivative(r rune) Regex {
	return NewComp(c.r.Derivative(r))
}

// Accepting returns true if the complemented regex
// does not accepting.
func (c *comp) Accepting() bool {
	return !c.r.Accepting()
}

type kleene struct {
	r Regex
}

// NewKleene creates a regex that accepts the Kleene star of a regex.
func NewKleene(r Regex) Regex {
	return &kleene{
		r: r,
	}
}

func (k *kleene) String() string {
	return fmt.Sprintf("(%v)*", k.r)
}

// Derivative returns the concatenation of the derivative
// of the Kleene star'd regex and the regex.
func (k *kleene) Derivative(r rune) Regex {
	return NewConcat(k.r.Derivative(r), k)
}

// Accepting returns true.
func (k *kleene) Accepting() bool {
	return true
}

var (
	_ Regex = &empty{}
	_ Regex = &epsilon{}
	_ Regex = &any{}
	_ Regex = &char{}
	_ Regex = &union{}
	_ Regex = &comp{}
	_ Regex = &kleene{}
)
