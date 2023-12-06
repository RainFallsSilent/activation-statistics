package common

type SetCount map[string]int

func NewSetCount() SetCount {
	s := make(map[string]int)
	return s
}

func (s SetCount) Add(item string) {
	s[item] += 1
}

func (s SetCount) Remove(item string) {
	delete(s, item)
}

func (s SetCount) Contains(item string) bool {
	return s[item] > 0
}

func (s SetCount) Size() int {
	return len(s)
}

func (s SetCount) Values() []string {
	values := make([]string, 0, len(s))
	for item := range s {
		values = append(values, item)
	}
	return values
}
