package cpu

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hculpan/kabbit/pkg/executable"
	"github.com/hculpan/kabbit/pkg/opcodes"
)

type MonitorFunc func(cpu *Cpu, lastError *error)

type Cpu struct {
	StackPointer       int
	InstructionPointer int
	Stack              []int32
	Code               []int32
	Heap               []int32

	Monitor MonitorFunc

	halted    bool
	stackSize int
	heapSize  int
	codeSize  int
}

func NewCpu(file *executable.ExecutableFile, monitorFunc MonitorFunc) *Cpu {
	return &Cpu{
		StackPointer:       0,
		InstructionPointer: 0,
		Stack:              make([]int32, file.Header.StackSize),
		Code:               file.Code,
		Heap:               file.Data,
		halted:             false,
		stackSize:          int(file.Header.StackSize),
		codeSize:           len(file.Code),
		heapSize:           int(file.Header.HeapSize),
		Monitor:            monitorFunc,
	}
}

func (c *Cpu) Run() error {
	c.halted = false

	if c.Monitor != nil {
		c.Monitor(c, nil)
	}

	for !c.halted {
		err := c.Step()
		if c.Monitor != nil {
			c.Monitor(c, &err)
		}
		if err != nil {
			c.halted = true
			return err
		}
	}

	c.halted = true
	return nil
}

func (c *Cpu) IsHalted() bool {
	return c.halted
}

func (c *Cpu) Step() error {
	opcode := c.Code[c.InstructionPointer]
	param := c.Code[c.InstructionPointer+1]
	switch opcode {
	case opcodes.PUSH:
		if err := c.push(param); err != nil {
			return err
		}
	case opcodes.POP:
		if _, err := c.pop(); err != nil {
			return err
		}
	case opcodes.DUP:
		v, err := c.pop()
		if err != nil {
			return err
		}

		err = c.push(v)
		if err != nil {
			return err
		}

		err = c.push(v)
		if err != nil {
			return err
		}
	case opcodes.OUT:
		v, err := c.pop()
		if err != nil {
			return err
		}

		fmt.Println(v)
	case opcodes.IN:
		var num int32 = 0
		if v, err := readInteger("-> "); err == nil {
			num = int32(v)
		}
		if err := c.push(num); err != nil {
			return err
		}
	case opcodes.ST:
		if param >= int32(c.heapSize) {
			return errors.New(fmt.Sprintf("invalid memory location %d", param))
		}

		v, err := c.pop()
		if err != nil {
			return err
		}

		c.Heap[param] = v
	case opcodes.LD:
		if param >= int32(c.heapSize) {
			return errors.New(fmt.Sprintf("invalid memory location %d", param))
		}

		v := c.Heap[param]
		if err := c.push(v); err != nil {
			return err
		}
	case opcodes.ADD, opcodes.SUB, opcodes.MUL, opcodes.DIV, opcodes.AND, opcodes.OR, opcodes.XOR:
		v1, err := c.pop()
		if err != nil {
			return err
		}
		v2, err := c.pop()
		if err != nil {
			return err
		}

		total := c.binaryOp(opcode, v1, v2)

		if err := c.push(total); err != nil {
			return err
		}
	case opcodes.DEC:
		v, err := c.pop()
		if err != nil {
			return err
		}
		v--
		if err := c.push(v); err != nil {
			return err
		}
	case opcodes.INC:
		v, err := c.pop()
		if err != nil {
			return err
		}
		v++
		if err := c.push(v); err != nil {
			return err
		}
	case opcodes.JMP:
		if param < int32(c.codeSize) {
			c.InstructionPointer = int(param)
			return nil
		} else {
			c.halted = true
			return errors.New(fmt.Sprintf("invalid jmp address %d", param))
		}
	case opcodes.JIF:
		v, err := c.pop()
		if err != nil {
			return err
		}

		if v != 0 {
			if param < int32(c.codeSize) {
				c.InstructionPointer = int(param)
				return nil
			} else {
				c.halted = true
				return errors.New(fmt.Sprintf("invalid jmp address %d", param))
			}
		}
	case opcodes.HALT:
		c.halted = true
	default:
		c.halted = true
		return errors.New("invalid instruction")
	}

	c.InstructionPointer += 2

	return nil
}

func (c *Cpu) binaryOp(opcode int32, v1 int32, v2 int32) int32 {
	switch opcode {
	case opcodes.ADD:
		return v1 + v2
	case opcodes.SUB:
		return v1 - v2
	case opcodes.MUL:
		return v1 * v2
	case opcodes.DIV:
		return v1 / v2
	case opcodes.AND:
		if v1 != 0 && v2 != 0 {
			return 1
		} else {
			return 0
		}
	case opcodes.OR:
		if v1 != 0 || v2 != 0 {
			return 1
		} else {
			return 0
		}
	case opcodes.XOR:
		if (v1 != 0 || v2 != 0) && !(v1 != 0 && v2 != 0) {
			return 1
		} else {
			return 0
		}
	default:
		c.halted = true
		return 0
	}

}

func (c *Cpu) pop() (int32, error) {
	if c.StackPointer < 1 {
		return 0, errors.New("stack underflow")
	}

	c.StackPointer--
	result := c.Stack[c.StackPointer]

	return result, nil
}

func (c *Cpu) push(v int32) error {
	if c.StackPointer >= c.stackSize {
		return errors.New("stack overflow")
	}

	c.Stack[c.StackPointer] = v
	c.StackPointer++

	return nil
}

func readInteger(prompt string) (int32, error) {
	reader := bufio.NewReader(os.Stdin)
	var number int32
	for {
		fmt.Print(prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			return 0, err
		}

		input = strings.TrimSpace(input)
		n, err := strconv.Atoi(input)
		if err == nil {
			number = int32(n)
			break
		}
	}

	return number, nil
}
