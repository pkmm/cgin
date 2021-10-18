NAME=cgin
BINDIR=bin
GOBUILD=CGO_ENABLED=0 go build -idflags '-w -s -buildid='

all: linux win32 win64

linux:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

win64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe

win32:
	GOARCH=386 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe

test: test-linux test-win64 test-win32

test-linux:
	GOARCH=amd64 GOOS=linux go test

test-win64:
	GOARCH=amd64 GOOS=windows go test

test-win32:
	GOARCH=386 GOOS=windows go test

clean:
	rm $(BINDIR)/*