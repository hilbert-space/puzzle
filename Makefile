all: puzzle-c puzzle-go

puzzle-c: main.c
	$(CC) $< -lpthread -o $@

puzzle-go: main.go
	go build -o $@ $<

clean:
	$(RM) puzzle-c puzzle-go
