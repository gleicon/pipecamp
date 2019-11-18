all: cli

export GO111MODULE=on

cli:
	go build -v -o pipecamp 

clean:
	rm -f pipecamp

test:
	go test -v
	
