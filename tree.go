package dr

type regexTree struct {
	stack []Regex
}

func (t *regexTree) get() Regex {
	if len(t.stack) != 1 {
		panic("stack length not 1")
	}
	return t.stack[0]
}

func (t *regexTree) push(r Regex) {
	t.stack = append(t.stack, r)
}

func (t *regexTree) pop() Regex {
	var r Regex
	r, t.stack = t.stack[len(t.stack)-1], t.stack[:len(t.stack)-1]
	return r
}

func (t *regexTree) empty() bool {
	return len(t.stack) == 0
}

func (t *regexTree) char(r rune) {
	t.push(NewChar(r))
}

func (t *regexTree) any() {
	t.push(NewAny())
}

func (t *regexTree) kleene() {
	r := t.pop()
	t.push(NewKleene(r))
}

func (t *regexTree) comp() {
	r := t.pop()
	t.push(NewComp(r))
}

func (t *regexTree) concat() {
	a := t.pop()
	b := t.pop()

	t.push(NewConcat(b, a))
}

func (t *regexTree) union() {
	a := t.pop()
	b := t.pop()

	t.push(NewUnion(b, a))
}
