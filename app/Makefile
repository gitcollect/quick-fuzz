PKGS = quickio

CFLAGS = \
	-O3 \
	-fPIC \
	-shared \
	$(shell pkg-config --cflags $(PKGS))

LDFLAGS = \
	$(shell pkg-config --libs $(PKGS))

LIB_DIR = $(DESTDIR)/usr/lib/quickio

all: quick-fuzz.so

quick-fuzz.so: quick-fuzz.c
	$(CC) $(CFLAGS) $< -o $@ $(LDFLAGS)

install: quick-fuzz.so
	mkdir -p $(LIB_DIR)
	install -m 644 quick-fuzz.so $(LIB_DIR)

uninstall:
	rm -f $(LIB_DIR)/quick-fuzz.so

clean:
	rm -f quick-fuzz.so
