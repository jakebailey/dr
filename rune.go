package dr

func firstRune(s string) rune {
	var r rune
	for _, c := range s {
		r = c
		break
	}
	return r
}

func lastRune(s string) rune {
	var r rune
	for _, c := range s {
		r = c
	}
	return r
}
