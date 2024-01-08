package opcodes

import (
	"errors"
	"fmt"
)

type OperandType int

const (
	NONE OperandType = iota
	INT32
)

const (
	PUSH  = 1
	POP   = 2
	ADD   = 3
	SUB   = 4
	MUL   = 5
	DIV   = 6
	DUP   = 7
	DEC   = 8
	INC   = 9
	JMP   = 10
	JIF   = 11
	OUT   = 20
	IN    = 21
	ST    = 30
	LD    = 31
	AND   = 40
	OR    = 41
	XOR   = 42
	ISEQ  = 50
	ISGT  = 51
	ISGTE = 52
	HALT  = 0xFFFF
	WD    = 0
)

type Instruction struct {
	Pneumonic string
	Opcode    uint32
	Param     OperandType
	Dataop    bool
}

var opcodes map[string]Instruction = map[string]Instruction{
	"invalid": {Pneumonic: "invalid", Opcode: 0, Param: NONE},
	"push":    {Pneumonic: "push", Opcode: 1, Param: INT32},
	"pop":     {Pneumonic: "pop", Opcode: 2, Param: NONE},
	"add":     {Pneumonic: "add", Opcode: 3, Param: NONE},
	"sub":     {Pneumonic: "sub", Opcode: 4, Param: NONE},
	"mul":     {Pneumonic: "mul", Opcode: 5, Param: NONE},
	"div":     {Pneumonic: "div", Opcode: 6, Param: NONE},
	"dup":     {Pneumonic: "dup", Opcode: 7, Param: NONE},
	"dec":     {Pneumonic: "dec", Opcode: 8, Param: NONE},
	"inc":     {Pneumonic: "inc", Opcode: 9, Param: NONE},
	"jmp":     {Pneumonic: "jmp", Opcode: 10, Param: INT32},
	"jif":     {Pneumonic: "jif", Opcode: 11, Param: INT32},
	"out":     {Pneumonic: "out", Opcode: 20, Param: NONE},
	"in":      {Pneumonic: "in", Opcode: 21, Param: NONE},
	"st":      {Pneumonic: "st", Opcode: 30, Param: INT32},
	"ld":      {Pneumonic: "ld", Opcode: 31, Param: INT32},
	"and":     {Pneumonic: "and", Opcode: 40, Param: NONE},
	"or":      {Pneumonic: "or", Opcode: 41, Param: NONE},
	"xor":     {Pneumonic: "xor", Opcode: 42, Param: NONE},
	"iseq":    {Pneumonic: "iseq", Opcode: 50, Param: NONE},
	"isgt":    {Pneumonic: "isgt", Opcode: 51, Param: NONE},
	"isgte":   {Pneumonic: "isgte", Opcode: 52, Param: NONE},
	"halt":    {Pneumonic: "halt", Opcode: 0xFFFF, Param: NONE},
	"wd":      {Pneumonic: "wd", Opcode: 0, Param: INT32, Dataop: true},
}

func GetPneumonic(opcode uint32) string {
	i, err := GetInstructionByOpcode(opcode)
	if err != nil {
		return "unkn"
	}

	return i.Pneumonic
}

func GetInstructionByOpcode(opcode uint32) (*Instruction, error) {
	for _, v := range opcodes {
		if v.Opcode == uint32(opcode) {
			return &v, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("unrecognized opcode %d", opcode))
}

func GetInstructionByPneumonic(pneumonic string) (*Instruction, error) {
	for k, v := range opcodes {
		if pneumonic == k {
			return &v, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("unsupported operation '%s'", pneumonic))
}
