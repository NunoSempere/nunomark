dev:
	go run main.go README.md | tee README.html

view:
	rosenrot file://$$(realpath ./README.html)

build:
	go build -ldflags "-s -w" -o nunomark main.go

install:
	cp nunomark /usr/bin/nunomark

