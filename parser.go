package dr

func Parse(s string) (Regex, error) {
	p := parseTree{Buffer: s}
	p.Init()

	if err := p.Parse(); err != nil {
		return nil, err
	}

	p.Execute()
	return p.get(), nil
}

func MustParse(s string) Regex {
	r, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return r
}

//go:generate peg -inline -switch peg.peg
