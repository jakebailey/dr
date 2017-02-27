package dr

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleRoot
	ruleRegex
	ruleUnion
	ruleConcat
	ruleUnary
	ruleComp
	ruleKleene
	ruleFactor
	ruleChar
	ruleAction0
	ruleAction1
	ruleAny
	ruleAction2
	ruleAction3
	ruleAction4
	rulePegText
	ruleAction5
	ruleAction6
	ruleAction7
)

var rul3s = [...]string{
	"Unknown",
	"Root",
	"Regex",
	"Union",
	"Concat",
	"Unary",
	"Comp",
	"Kleene",
	"Factor",
	"Char",
	"Action0",
	"Action1",
	"Any",
	"Action2",
	"Action3",
	"Action4",
	"PegText",
	"Action5",
	"Action6",
	"Action7",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Printf(" ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Printf("%v %v\n", rule, quote)
			} else {
				fmt.Printf("\x1B[34m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(buffer string) {
	node.print(false, buffer)
}

func (node *node32) PrettyPrint(buffer string) {
	node.print(true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	if tree := t.tree; int(index) >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	t.tree[index] = token32{
		pegRule: rule,
		begin:   begin,
		end:     end,
	}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type parseTree struct {
	regexTree

	Buffer string
	buffer []rune
	rules  [20]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *parseTree) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *parseTree) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *parseTree
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *parseTree) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *parseTree) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for _, token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:
			p.union()
		case ruleAction1:
			p.concat()
		case ruleAction2:

		case ruleAction3:
			p.comp()
		case ruleAction4:
			p.kleene()
		case ruleAction5:
			p.char(firstRune(text))
		case ruleAction6:
			p.char(lastRune(text))
		case ruleAction7:
			p.any()

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *parseTree) Init() {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 Root <- <(Regex !.)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
				if !_rules[ruleRegex]() {
					goto l0
				}
				{
					position2, tokenIndex2 := position, tokenIndex
					if !matchDot() {
						goto l2
					}
					goto l0
				l2:
					position, tokenIndex = position2, tokenIndex2
				}
				add(ruleRoot, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 Regex <- <Union> */
		func() bool {
			position3, tokenIndex3 := position, tokenIndex
			{
				position4 := position
				if !_rules[ruleUnion]() {
					goto l3
				}
				add(ruleRegex, position4)
			}
			return true
		l3:
			position, tokenIndex = position3, tokenIndex3
			return false
		},
		/* 2 Union <- <((Concat !('+' Union)) / (Concat '+' Union Action0))> */
		func() bool {
			position5, tokenIndex5 := position, tokenIndex
			{
				position6 := position
				{
					position7, tokenIndex7 := position, tokenIndex
					if !_rules[ruleConcat]() {
						goto l8
					}
					{
						position9, tokenIndex9 := position, tokenIndex
						if buffer[position] != rune('+') {
							goto l9
						}
						position++
						if !_rules[ruleUnion]() {
							goto l9
						}
						goto l8
					l9:
						position, tokenIndex = position9, tokenIndex9
					}
					goto l7
				l8:
					position, tokenIndex = position7, tokenIndex7
					if !_rules[ruleConcat]() {
						goto l5
					}
					if buffer[position] != rune('+') {
						goto l5
					}
					position++
					if !_rules[ruleUnion]() {
						goto l5
					}
					{
						add(ruleAction0, position)
					}
				}
			l7:
				add(ruleUnion, position6)
			}
			return true
		l5:
			position, tokenIndex = position5, tokenIndex5
			return false
		},
		/* 3 Concat <- <((Unary !Concat) / (Unary Concat Action1 &Any Action2))> */
		func() bool {
			position11, tokenIndex11 := position, tokenIndex
			{
				position12 := position
				{
					position13, tokenIndex13 := position, tokenIndex
					if !_rules[ruleUnary]() {
						goto l14
					}
					{
						position15, tokenIndex15 := position, tokenIndex
						if !_rules[ruleConcat]() {
							goto l15
						}
						goto l14
					l15:
						position, tokenIndex = position15, tokenIndex15
					}
					goto l13
				l14:
					position, tokenIndex = position13, tokenIndex13
					if !_rules[ruleUnary]() {
						goto l11
					}
					if !_rules[ruleConcat]() {
						goto l11
					}
					{
						add(ruleAction1, position)
					}
					{
						position17, tokenIndex17 := position, tokenIndex
						{
							position18 := position
							add(ruleAny, position18)
						}
						position, tokenIndex = position17, tokenIndex17
					}
					{
						add(ruleAction2, position)
					}
				}
			l13:
				add(ruleConcat, position12)
			}
			return true
		l11:
			position, tokenIndex = position11, tokenIndex11
			return false
		},
		/* 4 Unary <- <(Comp / Kleene / Factor)> */
		func() bool {
			position20, tokenIndex20 := position, tokenIndex
			{
				position21 := position
				{
					position22, tokenIndex22 := position, tokenIndex
					{
						position24 := position
						if buffer[position] != rune('!') {
							goto l23
						}
						position++
						if !_rules[ruleFactor]() {
							goto l23
						}
						{
							add(ruleAction3, position)
						}
						add(ruleComp, position24)
					}
					goto l22
				l23:
					position, tokenIndex = position22, tokenIndex22
					{
						position27 := position
						if !_rules[ruleFactor]() {
							goto l26
						}
						if buffer[position] != rune('*') {
							goto l26
						}
						position++
						{
							add(ruleAction4, position)
						}
						add(ruleKleene, position27)
					}
					goto l22
				l26:
					position, tokenIndex = position22, tokenIndex22
					if !_rules[ruleFactor]() {
						goto l20
					}
				}
			l22:
				add(ruleUnary, position21)
			}
			return true
		l20:
			position, tokenIndex = position20, tokenIndex20
			return false
		},
		/* 5 Comp <- <('!' Factor Action3)> */
		nil,
		/* 6 Kleene <- <(Factor '*' Action4)> */
		nil,
		/* 7 Factor <- <(Char / ('(' Regex ')'))> */
		func() bool {
			position31, tokenIndex31 := position, tokenIndex
			{
				position32 := position
				{
					position33, tokenIndex33 := position, tokenIndex
					{
						position35 := position
						{
							position36, tokenIndex36 := position, tokenIndex
							{
								position38 := position
								{
									position39, tokenIndex39 := position, tokenIndex
									{
										switch buffer[position] {
										case '.':
											if buffer[position] != rune('.') {
												goto l39
											}
											position++
											break
										case ')':
											if buffer[position] != rune(')') {
												goto l39
											}
											position++
											break
										case '(':
											if buffer[position] != rune('(') {
												goto l39
											}
											position++
											break
										case '\\':
											if buffer[position] != rune('\\') {
												goto l39
											}
											position++
											break
										case '!':
											if buffer[position] != rune('!') {
												goto l39
											}
											position++
											break
										case '*':
											if buffer[position] != rune('*') {
												goto l39
											}
											position++
											break
										default:
											if buffer[position] != rune('+') {
												goto l39
											}
											position++
											break
										}
									}

									goto l37
								l39:
									position, tokenIndex = position39, tokenIndex39
								}
								if !matchDot() {
									goto l37
								}
								add(rulePegText, position38)
							}
							{
								add(ruleAction5, position)
							}
							goto l36
						l37:
							position, tokenIndex = position36, tokenIndex36
							if buffer[position] != rune('\\') {
								goto l42
							}
							position++
							{
								position43 := position
								{
									switch buffer[position] {
									case '.':
										if buffer[position] != rune('.') {
											goto l42
										}
										position++
										break
									case ')':
										if buffer[position] != rune(')') {
											goto l42
										}
										position++
										break
									case '(':
										if buffer[position] != rune('(') {
											goto l42
										}
										position++
										break
									case '\\':
										if buffer[position] != rune('\\') {
											goto l42
										}
										position++
										break
									case '!':
										if buffer[position] != rune('!') {
											goto l42
										}
										position++
										break
									case '*':
										if buffer[position] != rune('*') {
											goto l42
										}
										position++
										break
									default:
										if buffer[position] != rune('+') {
											goto l42
										}
										position++
										break
									}
								}

								add(rulePegText, position43)
							}
							{
								add(ruleAction6, position)
							}
							goto l36
						l42:
							position, tokenIndex = position36, tokenIndex36
							if buffer[position] != rune('.') {
								goto l34
							}
							position++
							{
								add(ruleAction7, position)
							}
						}
					l36:
						add(ruleChar, position35)
					}
					goto l33
				l34:
					position, tokenIndex = position33, tokenIndex33
					if buffer[position] != rune('(') {
						goto l31
					}
					position++
					if !_rules[ruleRegex]() {
						goto l31
					}
					if buffer[position] != rune(')') {
						goto l31
					}
					position++
				}
			l33:
				add(ruleFactor, position32)
			}
			return true
		l31:
			position, tokenIndex = position31, tokenIndex31
			return false
		},
		/* 8 Char <- <((<(!((&('.') '.') | (&(')') ')') | (&('(') '(') | (&('\\') '\\') | (&('!') '!') | (&('*') '*') | (&('+') '+')) .)> Action5) / ('\\' <((&('.') '.') | (&(')') ')') | (&('(') '(') | (&('\\') '\\') | (&('!') '!') | (&('*') '*') | (&('+') '+'))> Action6) / ('.' Action7))> */
		nil,
		/* 10 Action0 <- <{ p.union() }> */
		nil,
		/* 11 Action1 <- <{ p.concat() }> */
		nil,
		nil,
		/* 13 Action2 <- <{}> */
		nil,
		/* 14 Action3 <- <{ p.comp() }> */
		nil,
		/* 15 Action4 <- <{ p.kleene() }> */
		nil,
		nil,
		/* 17 Action5 <- <{ p.char(firstRune(text)) }> */
		nil,
		/* 18 Action6 <- <{ p.char(lastRune(text)) }> */
		nil,
		/* 19 Action7 <- <{ p.any() }> */
		nil,
	}
	p.rules = _rules
}
