.PHONY: build test install release

build test install:
	$(MAKE) -C godocdown $@

release: build
	godocdown/godocdown $(HOME)/go/src/pkg/strings > example.markdown
	(cd godocdown && ./godocdown -signature > README.markdown) || false
	cp godocdown/README.markdown .
