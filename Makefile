PACKAGES=mc mc/protocol
OUTFILE=mc
MAINFILE=src/main.go

VARS=GOPATH=`pwd`

both: format test build

format:
	$(VARS) go fmt $(PACKAGES)

test:
	$(VARS) go test -i $(PACKAGES)
	$(VARS) go test $(PACKAGES)

build:
	$(VARS) go build -o $(OUTFILE) $(MAINFILE)

clean:
	$(VARS) go clean
	rm -rf pkg
	rm -f mct
