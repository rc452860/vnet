export PATH := $(GOPATH)/bin:$(PATH)
LDFLAGS := -s -w
# The -w and -s flags reduce binary sizes by excluding unnecessary symbols and debug info

all:
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/vnet_darwin_amd64 ./cmd/server/server.go
	env CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -ldflags "$(LDFLAGS)" -o bin/vnet_freebsd_386 ./cmd/server/server.go
	env CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/vnet_freebsd_amd64 ./cmd/server/server.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "$(LDFLAGS)" -o bin/vnet_linux_386 ./cmd/server/server.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/vnet_linux_amd64 ./cmd/server/server.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags "$(LDFLAGS)" -o bin/vnet_linux_arm ./cmd/server/server.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/vnet_linux_arm64 ./cmd/server/server.go
	env CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags "$(LDFLAGS)" -o bin/vnet_windows_386.exe ./cmd/server/server.go
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/vnet_windows_amd64.exe ./cmd/server/server.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -ldflags "$(LDFLAGS)" -o bin/vnet_linux_mips64 ./cmd/server/server.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -ldflags "$(LDFLAGS)" -o bin/vnet_linux_mips64le ./cmd/server/server.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=mips GOMIPS=softfloat go build -ldflags "$(LDFLAGS)" -o bin/vnet_linux_mips ./cmd/server/server.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -ldflags "$(LDFLAGS)" -o bin/vnet_linux_mipsle ./cmd/server/server.go

clean:
	rm $(BINDIR)/*