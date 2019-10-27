package brainfuck

type State struct {
	data []byte
	pos  int
}

func (b *State) incPos() {
	b.pos++
}

func (b *State) decPos() {
	b.pos--
}

func (b *State) inc() {
	b.data[b.pos]++
}

func (b *State) dec() {
	b.data[b.pos]--
}

func (b *State) get() byte {
	return b.data[b.pos]
}

func (b *State) put(v byte) {
	b.data[b.pos] = v
}
