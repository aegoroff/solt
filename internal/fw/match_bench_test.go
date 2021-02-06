package fw

import (
	"math/rand"
	"testing"
)

const size = 100000
const nomatch = "sgfdgdsfg-s934ns4s0"

func BenchmarkNewExactMatch_NoMatch_NoBloom(b *testing.B) {
	data := generateRandomStringSlice(size, 50)
	m := NewExactMatch(data, NewNoneFilter())
	for i := 0; i < b.N; i++ {
		m.Match(nomatch)
	}
	b.ReportAllocs()
}

func BenchmarkNewExactMatch_NoMatch_Bloom(b *testing.B) {
	data := generateRandomStringSlice(size, 50)
	n := uint(len(data))
	bloom := NewBloomFilter(n)
	m := NewExactMatch(data, bloom)
	for i := 0; i < b.N; i++ {
		m.Match(nomatch)
	}
	b.ReportAllocs()
}

func BenchmarkNewExactMatch_Match_NoBloom(b *testing.B) {
	data := generateRandomStringSlice(size, 50)
	m := NewExactMatch(data, NewNoneFilter())
	s := data[size/2]
	for i := 0; i < b.N; i++ {
		m.Match(s)
	}
	b.ReportAllocs()
}

func BenchmarkNewExactMatch_Match_Bloom(b *testing.B) {
	data := generateRandomStringSlice(size, 50)
	m := NewExactMatch(data, NewBloomFilter(uint(len(data))))
	s := data[size/2]
	for i := 0; i < b.N; i++ {
		m.Match(s)
	}
	b.ReportAllocs()
}

func generateRandomStringSlice(num int, length int) []string {
	result := make([]string, num)
	for i := 0; i < num; i++ {
		l := 1 + rand.Intn(length)
		s := randomString(l)
		result[i] = s
	}
	return result
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		ix := rand.Intn(len(letters))
		s[i] = letters[ix]
	}
	return string(s)
}
