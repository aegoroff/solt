package cmd

// StringHashSet defines strings hash set
type StringHashSet map[string]interface{}

// IntHashSet defines integers hash set
type IntHashSet map[int]interface{}

// Items gets all set's items
func (m *StringHashSet) Items() []string {
	keys := make([]string, 0, len(*m))
	for k := range *m {
		keys = append(keys, k)
	}
	return keys
}

// Contains gets whether a key is presented within the set
func (m *StringHashSet) Contains(key string) bool {
	_, ok := (*m)[key]
	return ok
}

// Add adds new item into the set
func (m *StringHashSet) Add(key string) {
	(*m)[key] = nil
}
