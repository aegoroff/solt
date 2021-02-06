package fw

import (
	"github.com/willf/bloom"
	"strings"
)

type noneFilter struct {
}

// NewNoneFilter creates filter that do nothing
func NewNoneFilter() MatchFilter {
	return &noneFilter{}
}

func (*noneFilter) Match(string) bool { return true }
func (*noneFilter) Append(string)     {}

type bloomFilter struct {
	filter    *bloom.BloomFilter
	decorator func(s string) string
}

// NewBloomFilter creates new Bloom filter instance
func NewBloomFilter(sz uint) MatchFilter {
	filter := bloom.New(16*sz*2, 6)

	return &bloomFilter{
		filter:    filter,
		decorator: strings.ToUpper,
	}
}

func (b *bloomFilter) Append(s string) {
	b.filter.AddString(b.decorator(s))
}

func (b *bloomFilter) estimate(n uint) float64 {
	return b.filter.EstimateFalsePositiveRate(n)
}

func (b *bloomFilter) Match(s string) bool {
	return b.filter.TestString(b.decorator(s))
}
