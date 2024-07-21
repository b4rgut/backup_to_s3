.PHONY: build

build:
	GOOS=windows GOARCH=amd64 go build -o bin/cloud_backuper/cloud_backuper.exe -ldflags "-s -w" -a -installsuffix cgo src/main.go
	@cp -f config.yml bin/cloud_backuper/