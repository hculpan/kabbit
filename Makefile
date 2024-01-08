all: build

build:
	go build -o dist/kaba cmd/assembler/*.go 
	go build -o dist/kabc cmd/compiler/*.go 
	go build -o dist/kabv cmd/vm/*.go 

clean:
	rm -rf dist