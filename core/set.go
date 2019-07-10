package core

type Set map[string]struct{}

func (s *Set) Add(key string) {
	(*s)[key] = struct{}{}
}

func (s *Set) Has(key string) bool {
	_, ok := (*s)[key]
	return ok
}

func (s *Set) Copy() *Set {
	copy := make(Set)
	m := *s
	for v := range m {
		copy[v] = struct{}{}
	}
	return &copy
}
