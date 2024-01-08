package assembler

import (
	"fmt"

	"github.com/hculpan/kabbit/pkg/opcodes"
)

func Parse(l *Lexer) ([]Node, error) {
	result := make([]Node, 0)

	result = append(result, &ProgramNode{
		LineNo: 0,
	})

	currToken := l.NextToken()
	for currToken != nil && currToken.Type != TokenTypeEOF {
		switch currToken.Type {
		case TokenTypeDirective:
			node := &DirectiveNode{
				Directive: currToken.Literal[1:],
				LineNo:    currToken.LineNo,
			}
			result = append(result, node)
			if err := expectedToken(l, TokenTypeEOL, TokenTypeComment); err != nil {
				return nil, err
			}
		case TokenTypeLabel:
			node := &LabelNode{
				Name:   currToken.Literal[:len(currToken.Literal)-1],
				LineNo: currToken.LineNo,
			}
			result = append(result, node)
			token := l.NextToken()
			if token.Type != TokenTypeEOL {
				l.PushToken(token)
			}
		case TokenTypeIdentifier:
			node, err := parseInstruction(l, currToken)
			if err != nil {
				return nil, err
			}
			result = append(result, node)
		}

		currToken = l.NextToken()
	}

	return result, nil
}

func expectedToken(l *Lexer, tType ...int) error {
	token := l.NextToken()
	if token == nil {
		return fmt.Errorf("[%d] expected token, found nil", token.LineNo)
	}

	for _, tt := range tType {
		if tt == token.Type {
			return nil
		}
	}

	return fmt.Errorf("[%d] unexpected token type %s found", token.LineNo, TokenTypeNames[token.Type])
}

func parseInstruction(l *Lexer, currToken *Token) (Node, error) {
	instr, err := opcodes.GetInstructionByPneumonic(currToken.Literal)
	if err != nil {
		return nil, err
	}

	var operand string

	nextToken := l.NextToken()
	if instr.Param == opcodes.INT32 && (nextToken.Type == TokenTypeIdentifier || nextToken.Type == TokenTypeNumber) {
		operand = nextToken.Literal
	} else if instr.Param == opcodes.NONE && nextToken.Type == TokenTypeEOL {
		l.PushToken(nextToken)
		operand = ""
	} else if instr.Param == opcodes.NONE && nextToken.Type != TokenTypeEOL && nextToken.Type != TokenTypeEOF {
		return nil, fmt.Errorf("[%d] unexpected operand, found '%s'", nextToken.LineNo, nextToken.Literal)
	} else if instr.Param == opcodes.INT32 && nextToken.Type == TokenTypeEOL && nextToken.Type != TokenTypeEOF {
		return nil, fmt.Errorf("[%d] expected operand, found none", nextToken.LineNo)
	} else {
		return nil, fmt.Errorf("[%d] expected identifier or number, found '%s'", nextToken.LineNo, nextToken.Literal)
	}

	if err := expectedToken(l, TokenTypeEOL); err != nil {
		return nil, err
	}

	if instr.Dataop {
		dataNode := &DataNode{
			DataType: currToken.Literal,
			Value:    operand,
			LineNo:   currToken.LineNo,
		}

		return dataNode, nil
	} else {
		instrNode := &InstructionNode{
			Pneumonic: currToken.Literal,
			LineNo:    currToken.LineNo,
			Operand:   operand,
		}

		return instrNode, nil
	}
}
