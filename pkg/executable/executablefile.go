package executable

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

var Endian binary.ByteOrder = binary.BigEndian

type ExecutableFile struct {
	Filename string
	Header   FileHeader
	Code     []int32
	Data     []int32
}

func NewDefaultExecutableFile(filename string) *ExecutableFile {
	return &ExecutableFile{
		Filename: filename,
		Header: FileHeader{
			CodeSize:  0,
			StackSize: 1024 * 1024,
			HeapSize:  0,
		},
		Code: []int32{},
		Data: []int32{},
	}
}

func NewExecutableFile(filename string, header *FileHeader, code, data []int32) *ExecutableFile {
	return &ExecutableFile{
		Filename: filename,
		Header:   *header,
		Code:     code,
		Data:     data,
	}
}

func (e *ExecutableFile) SaveFile() error {
	f, err := os.Create(e.Filename)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to create file '%s'", e.Filename))
	}
	defer f.Close()

	e.Header.CodeSize = uint32(len(e.Code))
	e.Header.HeapSize = uint32(len(e.Data))

	err = binary.Write(f, Endian, e.Header)
	if err != nil {
		return err
	}

	err = binary.Write(f, Endian, e.Code)
	if err != nil {
		return err
	}

	err = binary.Write(f, Endian, e.Data)
	if err != nil {
		return err
	}

	return nil
}

func NewExecutableFromFile(filename string) (*ExecutableFile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the first 3 uint32 values (CodeSize, StackSize, HeapSize)
	header := FileHeader{}
	if err := binary.Read(file, Endian, &header); err != nil {
		return nil, err
	}

	// Read the remaining ints into the Code slice
	var code []int32 = []int32{}
	var data []int32 = []int32{}
	for i := 0; i < int(header.CodeSize); i++ {
		var value int32
		err := binary.Read(file, Endian, &value)
		if err != nil {
			break // Reached the end of the file
		}
		code = append(code, value)
	}

	for i := 0; i < int(header.HeapSize); i++ {
		var value int32
		err := binary.Read(file, Endian, &value)
		if err != nil {
			break // Reached the end of the file
		}
		data = append(data, value)
	}

	return &ExecutableFile{
		Filename: filename,
		Header:   header,
		Code:     code,
		Data:     data,
	}, nil
}
