package dr

import "testing"

func BenchmarkParseSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse("abc*+aad")
	}
}

func BenchmarkParseComplex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse("!(a+b*(asd(!d)))+(def)*")
	}
}

func BenchmarkMatchSimple(b *testing.B) {
	r := MustParse("abc*+aad")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Match(r, "abccccccccc")
	}
}

func BenchmarkMatchComplex(b *testing.B) {
	r := MustParse("!(a+b*(asd(!d)))+(def)*")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Match(r, "abccccccccc")
	}
}
