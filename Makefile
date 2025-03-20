.PHONY: all,linux,windows

all: linux windows

linux:
	@echo "Building for Linux"
	@GOOS=linux GOARCH=amd64 go build -C ./cli -o ../nbtreader

windows:
	@echo "Building for Windows"
	@GOOS=windows GOARCH=amd64 go build -C ./cli -o ../nbtreader.exe
