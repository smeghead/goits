VERSION = 0.0.01
COPYRIGHT = Copyright\ smeghead\ 2013\ -\ 2013
CC = 6g
LINKER = 6l

PROGRAM = visisonor

$(PROGRAM): data.6 handler.6 main.6
	$(LINKER) -o $@ main.6

handler.6: handler/*.go
	$(CC) -o $@ $^

data.6: data/*.go
	$(CC) -o $@ $^

main.6: *.go
	$(CC) -o $@ $^

.PHONY: clean run
clean:
	rm -f $(PROGRAM) *.6

run: $(PROGRAM)
	./$(PROGRAM)
