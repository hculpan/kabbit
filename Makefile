all: build

build:
	go build -o dist/kabbitc cmd/compiler/*.go 
	go build -o dist/kabbit cmd/vm/*.go 

clean:
	rm -rf dist