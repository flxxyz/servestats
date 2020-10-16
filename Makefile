ARCH_NAME=amd64
EXT_NAME=
ifeq ($(OS),Windows_NT)
	EXT_NAME=.exe
	OS_NAME=windows
	ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
		ARCH_NAME=amd64
	endif
	ifeq ($(PROCESSOR_ARCHITECTURE),x86)
		ARCH_NAME=386
	endif
else
	UNAME_S=$(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		OS_NAME=linux
	endif
	ifeq ($(UNAME_S),Darwin)
		OS_NAME=darwin
	endif

	UNAME_P=$(shell uname -m)
	ifeq ($(UNAME_P),x86_64)
		ARCH_NAME=amd64
	endif
	ifneq ($(filter %86,$(UNAME_P)),)
		ARCH_NAME=386
	endif
	ifneq ($(filter arm%,$(UNAME_P)),)
		ARCH_NAME=arm
	endif
	ifneq ($(filter aarch%,$(UNAME_P)),)
		ARCH_NAME=arm64
	endif
endif

FILE_NAME=ServerStatus_$(OS_NAME)_$(ARCH_NAME)$(EXT_NAME)

all: build

build: clean
	CGO_ENABLED=0 GOOS=$(OS_NAME) GOARCH=$(ARCH_NAME) go build -v -gcflags "all=-N -l" -trimpath -o $(FILE_NAME)

win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -v -gcflags "all=-N -l" -trimpath -o ServerStatus_windows_amd64.exe

server: build
	./$(FILE_NAME) server

client: build
	./$(FILE_NAME) client -t 1s -id c91f4a1b-865b-435e-8e0a-f35a0557e19d

system: build
	./$(FILE_NAME) system

traffic: build
	./$(FILE_NAME) traffic

clean:
	rm -rf $(FILE_NAME)
