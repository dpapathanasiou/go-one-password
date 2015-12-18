SHELL = /bin/sh
CC    = gcc

%: %.go
	go build -o $@ $^

TARGETS     = $(shell echo *.go | sed -e 's/\.go//g')
GUI_VERSION = $(shell echo *ui.go | sed -e 's/\.go//g')
CLI_VERSION = $(shell echo *cl.go | sed -e 's/\.go//g')

all: $(TARGETS)

gui: $(GUI_VERSION)

cli: $(CLI_VERSION)

clean:
	rm -f $(TARGETS)
