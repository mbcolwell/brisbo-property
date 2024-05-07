MAIN = brisboproperty
BIN = bin/
CMD = cmd/
PIPE = 
ARGS = --n 10 --m 20


all: build run

build:
	go build -o $(BIN)$(MAIN) $(CMD)$(MAIN).go

run:
	$(PIPE) ./$(BIN)$(MAIN) $(ARGS)

clean:
	rm -f $(BIN)$(MAIN) output.log
