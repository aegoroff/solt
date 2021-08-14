package fw

import (
	"github.com/willf/bloom"
	"strings"
)

type bloomFilter struct {
	filter    *bloom.BloomFilter
	decorator func(s string) string
}

// NewBloomFilter creates new Bloom filter instance
func NewBloomFilter(matches []string) Matcher {
	filter := bloom.New(16*uint(len(matches))*2, 6)

	d := strings.ToUpper
	for _, match := range matches {
		filter.AddString(d(match))
	}
	return &bloomFilter{
		filter:    filter,
		decorator: d,
	}
}

func (b *bloomFilter) estimate(n uint) float64 {
	return b.filter.EstimateFalsePositiveRate(n)
}

func (b *bloomFilter) Match(s string) bool {
	return b.filter.TestString(b.decorator(s))
}
