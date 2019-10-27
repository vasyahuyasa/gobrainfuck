package brainfuck

import (
	"errors"
	"fmt"
	"io"
)

var (
	ErrUnknownChar      = errors.New("unknown symbol")
	ErrUnknownOperation = errors.New("unknown operation")
	ErrNoLoopEnd        = errors.New("no loop end")
)

const (
	Unknown Operation = iota
	OpNext
	OpPrev
	OpInc
	OpDec
	OpPut
	OpGet
	OpLoop
	OpEndLoop

	DefaultSize = 30000
)

type Operation byte

type Interpreter struct {
	state *State

	// program counter
	pc      int
	pcstack stack

	// programm code
	tape []Operation

	// last interpretation error
	lastErr error
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		state: &State{
			pos:  0,
			data: make([]byte, DefaultSize),
		},
		pc:      0,
		pcstack: stack{},
		tape:    []Operation{},
	}
}

func (i *Interpreter) ParseString(str string) error {
	tape := make([]Operation, 0, len(str))
	for _, b := range []byte(str) {
		op, err := char2op(b)
		if err == ErrUnknownChar {
			continue
		}
		if err != nil {
			return err
		}
		tape = append(tape, op)
	}
	i.tape = tape

	return nil
}

func (i *Interpreter) Run() error {
	for i.next() {
	}

	if i.lastErr == io.EOF {
		return nil
	}

	return i.lastErr
}

func (i *Interpreter) next() bool {
	if i.pc >= len(i.tape) {
		i.lastErr = io.EOF
		return false
	}

	op := i.tape[i.pc]

	switch op {
	case OpNext:
		i.state.incPos()
		i.pc++
	case OpPrev:
		i.state.decPos()
		i.pc++
	case OpInc:
		i.state.inc()
		i.pc++
	case OpDec:
		i.state.dec()
		i.pc++
	case OpPut:
		v := i.state.get()
		_, err := fmt.Printf("%s", string(v))
		if err != nil {
			i.lastErr = err
			return false
		}
		i.pc++
	case OpGet:
		var s string
		for len(s) == 0 {
			fmt.Scanf("%s", &s)
		}
		i.state.put(s[0])
		i.pc++
	case OpLoop:
		v := i.state.get()
		if v != 0 {
			i.pcstack.push(i.pc)
			i.pc++
			break
		}

		loopEnd, err := i.findloopEnd(i.pc)
		if err != nil {
			i.lastErr = err
			return false
		}
		i.pc = loopEnd + 1
	case OpEndLoop:
		loopStart := i.pcstack.pop()
		v := i.state.get()
		if v == 0 {
			i.pc++
			break
		}

		i.pc = loopStart

	case Unknown:
		fallthrough
	default:
		i.lastErr = ErrUnknownOperation
		return false
	}

	return true
}

func (i *Interpreter) findloopEnd(from int) (int, error) {
	level := 0
	to := from
	for _, op := range i.tape[from:] {
		switch op {
		case OpLoop:
			level++
		case OpEndLoop:
			level--
			if level == 0 {
				return to, nil
			}
		}
		to++
	}

	return 0, ErrNoLoopEnd
}

func char2op(c byte) (Operation, error) {
	switch c {
	case '>':
		return OpNext, nil
	case '<':
		return OpPrev, nil
	case '+':
		return OpInc, nil
	case '-':
		return OpDec, nil
	case '.':
		return OpPut, nil
	case ',':
		return OpGet, nil
	case '[':
		return OpLoop, nil
	case ']':
		return OpEndLoop, nil
	default:
		return Unknown, ErrUnknownChar
	}
}
