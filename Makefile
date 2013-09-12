PACKAGES=mc mc/protocol mc/simulator nbt smpm httphelpers github.com/jeffh/goexpect
FMT_PACKAGES=$(PACKAGES)
OUTFILE=mc
MAINFILE=src/main.go

VARS=GOPATH=`pwd`

all: format test build

format:
	$(VARS) go fmt $(FMT_PACKAGES)

test:
	$(VARS) go test -i $(PACKAGES)
	$(VARS) go test $(PACKAGES)

build:
	$(VARS) go build -o $(OUTFILE) $(MAINFILE)

clean:
	$(VARS) go clean
	rm -rf pkg
	rm -f $(OUTFILE)
