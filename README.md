# Kabbit
A programming language and virtual machine written in Go.

The name is a mashup of two words: Kabluey and "bit" for binary digit, and I liked that it rhymes with "rabbit", implying fast (which this is most assuredly not, but one doesn't advertise that ;P).

This consists of 3 components: The assembler, the compiler, and the VM itself, of course. My current thinking is that the compiler will produce assembler code, which will then be processed by the assembler itself, thus making the assembly language my IL. Of course, I'm just doing this as a toy project with very little forethought or design, so this will probably end up changing.

# Assembler
```
Usage:
  assembler <input file> [flags]

Flags:
  -h, --help            help for assembler
  -o, --output string   Output file
```

The assembler takes an assembly code file as input (.kba extension), and produces an executable program (.kbx). By default, if the -o flag is not included, the executable will have the same name as the input file, but with the "kbx" extension.

## How to build
``` make build ``` or ``` make all ```

Executable files will be created in the ```dist``` directory.

## Clean
``` make clean ```

Removes the ```dist``` directory.
