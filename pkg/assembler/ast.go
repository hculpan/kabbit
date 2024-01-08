package assembler

import "fmt"

type Node interface {
	GetLineNo() int
	GetDescription() string
}

type ProgramNode struct {
	LineNo int
}

func (n *ProgramNode) GetLineNo() int {
	return n.LineNo
}

func (n *ProgramNode) GetDescription() string {
	return "ProgramNode"
}

type DirectiveNode struct {
	LineNo    int
	Directive string
}

func (n *DirectiveNode) GetLineNo() int {
	return n.LineNo
}

func (n *DirectiveNode) GetDescription() string {
	return fmt.Sprintf("DirectiveNode [%s]", n.Directive)
}

type DataNode struct {
	LineNo   int
	DataType string
	Value    string
}

func (n *DataNode) GetLineNo() int {
	return n.LineNo
}

func (n *DataNode) GetDescription() string {
	return fmt.Sprintf("DataNode [%s  %s]", n.DataType, n.Value)
}

type InstructionNode struct {
	LineNo    int
	Pneumonic string
	Operand   string
}

func (n *InstructionNode) GetLineNo() int {
	return n.LineNo
}

func (n *InstructionNode) GetDescription() string {
	return fmt.Sprintf("InstructionNode [%s  %s]", n.Pneumonic, n.Operand)
}

type LabelNode struct {
	LineNo int
	Name   string
}

func (n *LabelNode) GetLineNo() int {
	return n.LineNo
}

func (n *LabelNode) GetDescription() string {
	return fmt.Sprintf("LabelNode [%s]", n.Name)
}
