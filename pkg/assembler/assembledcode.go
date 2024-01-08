package assembler

import "github.com/hculpan/kabbit/pkg/executable"

type AssembledCode struct {
	Data []int32
	Code []int32
}

func (a *AssembledCode) NewFileHeader() *executable.FileHeader {
	return &executable.FileHeader{
		CodeSize:  uint32(len(a.Code)),
		StackSize: 1024 * 1024,
		HeapSize:  uint32(len(a.Data)),
	}
}
