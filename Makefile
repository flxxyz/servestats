all: build

build:
	go build -v -gcflags "all=-N -l" -trimpath -o ServerStatus

server: build
	./ServerStatus server

client: build
	./ServerStatus client -t 1s -id c91f4a1b-865b-435e-8e0a-f35a0557e19d
