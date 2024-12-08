build:
	go build -ldflags "-s -w" -o nunomark main.go

install:
	cp nunomark /usr/bin/nunomark
