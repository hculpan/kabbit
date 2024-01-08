package executable

type FileHeader struct {
	CodeSize  uint32 // in 4-byte
	StackSize uint32 // number of stack elements (64-bit)
	HeapSize  uint32 //
}
