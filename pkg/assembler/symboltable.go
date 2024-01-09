package assembler

import (
	"errors"
	"fmt"
)

var symbolTable map[string]int32

func init() {
	symbolTable = make(map[string]int32)
	symbolTable["index"] = 0
}

func AddSymbol(name string, value int32) error {
	_, ok := symbolTable[name]
	if ok {
		return errors.New(fmt.Sprintf("symbol '%s' already defined", name))
	}

	symbolTable[name] = value
	return nil
}

func ReplaceSymbol(name string, value int32) {
	symbolTable[name] = value
}

func GetSymbolValue(name string) (int32, error) {
	v, ok := symbolTable[name]
	if !ok {
		return -1, errors.New(fmt.Sprintf("symbol '%s' undefined", name))
	}

	return v, nil
}

func DisplaySymbolTable() string {
	result := "Symbol Table:\n"

	for k, v := range symbolTable {
		result += fmt.Sprintf("%15s : %08X [%d]\n", k, v, v)
	}

	return result
}
