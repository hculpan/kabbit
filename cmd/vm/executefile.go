package main

import (
	"fmt"
	"strings"

	"github.com/hculpan/kabbit/pkg/cpu"
	"github.com/hculpan/kabbit/pkg/executable"
	"github.com/hculpan/kabbit/pkg/opcodes"
)

func ExecuteFile(inputFile string, disassemble bool, trace bool) error {
	ef, err := executable.NewExecutableFromFile(inputFile)
	if err != nil {
		return err
	}

	if disassemble {
		return disassembleFile(ef)
	}

	var monitorFunc cpu.MonitorFunc = nil
	if trace {
		monitorFunc = Monitor
	}
	cpu := cpu.NewCpu(ef, monitorFunc)
	return cpu.Run()
}

func decode(opcode, param int32) string {
	instr, _ := opcodes.GetInstructionByOpcode(uint32(opcode))

	if instr.Param == opcodes.NONE {
		return fmt.Sprintf("%8s       ", strings.ToUpper(instr.Pneumonic))
	}
	return fmt.Sprintf("%8s %6X", strings.ToUpper(instr.Pneumonic), param)
}

func Monitor(c *cpu.Cpu, lastError *error) {
	if c.IsHalted() {
		return
	}
	stack := ""

	max := 3
	for i := c.StackPointer - 1; i > -1; i-- {
		if i != 0 {
			stack += fmt.Sprintf("%08X, ", c.Stack[i])
		} else {
			stack += fmt.Sprintf("%08X", c.Stack[i])
		}
		max--
		if max == 0 {
			break
		}
	}

	mem := ""
	max = 3
	for i := 0; i < max && i < len(c.Heap); i++ {
		if i < max-1 {
			mem += fmt.Sprintf("%08X ", c.Heap[i])
		} else {
			mem += fmt.Sprintf("%08X", c.Heap[i])
		}
	}

	fmt.Printf("  IP:%08X    %-10s    SP:%08X  Stack: [%-28s]    Mem:(%s)\n",
		c.InstructionPointer, decode(c.Code[c.InstructionPointer], c.Code[c.InstructionPointer+1]), c.StackPointer, stack, mem)
}

func disassembleFile(ef *executable.ExecutableFile) error {
	fmt.Printf("Disassembling file '%s'", ef.Filename)

	fmt.Println("\nHeader:")
	fmt.Printf("  Code Size : %10d [%08X]\n", ef.Header.CodeSize, ef.Header.CodeSize)
	fmt.Printf("  Stack Size: %10d [%08X]\n", ef.Header.StackSize, ef.Header.StackSize)
	fmt.Printf("  Heap Size : %10d [%08X]\n", ef.Header.HeapSize, ef.Header.HeapSize)

	fmt.Println()

	fmt.Println("Code:")
	for i := 0; i < len(ef.Code); i += 2 {
		instr, err := opcodes.GetInstructionByOpcode(uint32(ef.Code[i]))
		if err != nil {
			fmt.Println("  unkn")
		} else if instr.Param == opcodes.NONE {
			fmt.Printf("  %08X\t%08X:%08X\t%4s\n", i, ef.Code[i], ef.Code[i+1], opcodes.GetPneumonic(uint32(ef.Code[i])))
		} else {
			fmt.Printf("  %08X\t%08X:%08X\t%4s\t%08X [%d]\n", i, ef.Code[i], ef.Code[i+1], opcodes.GetPneumonic(uint32(ef.Code[i])), ef.Code[i+1], ef.Code[i+1])
		}
	}

	return nil
}
