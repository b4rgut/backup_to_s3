.PHONY: all

all: build copy_config

build:
	GOOS=windows GOARCH=amd64 go build -o bin/cloud_backuper/cloud_backuper.exe -ldflags "-s -w" -a -installsuffix cgo cmd/main.go

copy_config:
	@cp -f config.yml bin/cloud_backuper/
