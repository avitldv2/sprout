PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
GO ?= go

.PHONY: build install clean test

build:
	$(GO) build -o sprout ./cmd/sprout

install: build
	mkdir -p $(DESTDIR)$(BINDIR)
	cp sprout $(DESTDIR)$(BINDIR)/sprout
	chmod 755 $(DESTDIR)$(BINDIR)/sprout

clean:
	rm -f sprout

test:
	$(GO) test ./...

