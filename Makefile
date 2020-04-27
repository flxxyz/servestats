all: build

build:
	go build -v -gcflags "all=-N -l" -trimpath -o ServerStatus

server: build
	./ServerStatus server

client: build
	./ServerStatus client -h 127.0.0.1 -t 3s
