.PHONY: all clean distclean

all:	gh-release-check
clean:
distclean:	clean; go clean

gh-release-check:	$(wildcard *.go) $(wildcard */*.go)
	go build -o $@
