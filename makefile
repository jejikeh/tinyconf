CFLAGS=-Wall -Wextra -Wswitch-enum -std=c11 -pedantic
LIBS=

ambient: main.c
	$(CC) $(CFLAGS) -o build/ambient main.c $(LIBS)