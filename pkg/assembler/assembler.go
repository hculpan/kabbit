package assembler

import (
	"errors"
	"fmt"
	"os"
)

type Assembler struct {
	debugInfo bool
}

func NewAssembler(debugInfo bool) *Assembler {
	return &Assembler{
		debugInfo: debugInfo,
	}
}

func (a *Assembler) Assemble(input string) (*AssembledCode, error) {
	l := NewLexer(input)
	nodes, err := Parse(l)
	if err != nil {
		return nil, err
	}

	if len(nodes) == 0 {
		return nil, errors.New("empty node set")
	} else if _, ok := nodes[0].(*ProgramNode); !ok {
		return nil, fmt.Errorf("expected program node, found %s", nodes[0].GetDescription())
	}

	code, data, err := Generate(nodes, a.debugInfo)
	if err != nil {
		return nil, err
	}

	return &AssembledCode{
		Data: data,
		Code: code,
	}, nil
}

func (a *Assembler) AssembleFromFile(inputFile string) (*AssembledCode, error) {
	fileBytes, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Read %d bytes from file %s\n", len(fileBytes), inputFile)

	fileString := string(fileBytes)

	return a.Assemble(fileString)
}
