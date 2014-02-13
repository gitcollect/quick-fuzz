PKGS = quickio

CFLAGS = \
	-O3 \
	-fPIC \
	-shared \
	$(shell pkg-config --cflags $(PKGS))

LDFLAGS = \
	$(shell pkg-config --libs $(PKGS))

all: fuzzer.so fuzzer.ini

fuzzer.so: fuzzer.c
	$(CC) $(CFLAGS) $< -o $@ $(LDFLAGS)

fuzzer.ini: fuzzer.ini.in
	sed 's|{HERE}|$(CURDIR)|g' $< > $@

clean:
	rm -f fuzzer.so
	rm -f fuzzer.ini
