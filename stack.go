package brainfuck

type stack struct {
	data []int
}

func (s *stack) push(v int) {
	s.data = append(s.data, v)
}

func (s *stack) pop() int {

	last := len(s.data) - 1
	v := s.data[last]
	s.data = s.data[:last]
	return v
}
