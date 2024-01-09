package assembler

import (
	"fmt"
	"strconv"

	"github.com/hculpan/kabbit/pkg/opcodes"
)

func Generate(nodes []Node, debugInfo bool) ([]int32, []int32, error) {
	code := []int32{}
	data := []int32{0}

	if debugInfo {
		fmt.Println("\nAST:")
		for _, node := range nodes {
			fmt.Printf("  %08X\t%s\n", node.GetLineNo(), node.GetDescription())
		}
	}

	if err := pass1(nodes); err != nil { // find labels
		return nil, nil, err
	}

	if debugInfo {
		fmt.Println()
		fmt.Println(DisplaySymbolTable())
	}

	// pass 2 - generate code
	for _, node := range nodes {
		switch n := node.(type) {
		case *InstructionNode:
			instr, err := opcodes.GetInstructionByPneumonic(n.Pneumonic)
			if err != nil {
				return nil, nil, err
			}
			code = append(code, int32(instr.Opcode))
			if instr.Param == opcodes.INT32 {
				if v, err := getOperandValue(n.Operand, n.LineNo); err != nil {
					return nil, nil, err
				} else {
					code = append(code, v)
				}
			} else {
				code = append(code, 0)
			}
		case *DataNode:
			if v, err := getOperandValue(n.Value, n.LineNo); err != nil {
				return nil, nil, err
			} else {
				data = append(data, v)
			}
		}
	}

	return code, data, nil
}

func getOperandValue(operand string, lineNo int) (int32, error) {
	if v, err := strconv.Atoi(operand); err == nil {
		return int32(v), nil
	}

	// failed to convert, so assuming it's a symbol
	if v, err := GetSymbolValue(operand); err == nil {
		return v, nil
	}

	return 0, fmt.Errorf("[%d] unknown symbol '%s'", lineNo, operand)
}

func pass1(nodes []Node) error {
	codeLoc := 0
	dataLoc := 1
	for idx, node := range nodes {
		switch n := node.(type) {
		case *DataNode:
			dataLoc++
		case *InstructionNode:
			codeLoc += 2
		case *LabelNode:
			if addr, err := findNextNode(nodes, idx, codeLoc, dataLoc); err != nil {
				return err
			} else {
				AddSymbol(n.Name, int32(addr))
			}
		}
	}

	return nil
}

func findNextNode(nodes []Node, idx, codeLoc, dataLoc int) (int, error) {
	for i := idx + 1; i < len(nodes); i++ {
		switch nodes[i].(type) {
		case *DataNode:
			return dataLoc, nil
		case *InstructionNode:
			return codeLoc, nil
		}
	}

	return 0, fmt.Errorf("[%d] no valid nodes found for label", nodes[idx].GetLineNo())
}
