package main

type StringSet map[string]interface{}

func (set StringSet) Contains(value string) bool {
	_, ok := set[value]
	return ok
}

func (set StringSet) Add(value string) {
	set[value] = struct{}{}
}

func NewStringSet(items []string) StringSet {
	set := StringSet{}

	for _, item := range items {
		set.Add(item)
	}

	return set
}
