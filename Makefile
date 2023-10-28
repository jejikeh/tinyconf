EXAMPLE_FOLDER = examples
BINARY_FOLDER = $(EXAMPLE_FOLDER)/binary
DIS_FOLDER = $(EXAMPLE_FOLDER)/dis

run_fib:
	go run . -run -debug -i $(EXAMPLE_FOLDER)/fib.naive

run_fib_x:
	go run . -run -x -debug -i $(BINARY_FOLDER)/fib

build_fib:
	go run . -build -i $(EXAMPLE_FOLDER)/fib.naive -o $(BINARY_FOLDER)/fib -debug

dis_fib_o:
	go run . -dis -i $(BINARY_FOLDER)/fib -o $(DIS_FOLDER)/fib.naive

dis_fib:
	go run . -dis -i $(BINARY_FOLDER)/fib

lex_fib:
	go run . -lex -i $(EXAMPLE_FOLDER)/fib.naive

tests:
	go test ./...