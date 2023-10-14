CFLAGS=-Wall -Wextra -Wswitch-enum -std=c11 -pedantic
LIBS=

glow: main.c
	$(CC) $(CFLAGS) -o build/glow main.c $(LIBS)