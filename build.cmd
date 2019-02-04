
@echo off
choice /C wl /M "windows,linux"
if errorlevel 2 goto linux
if errorlevel 1 goto windows
:windows
echo build windows amd64
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-s -w" -o vnet.exe  .\cmd\server\server.go
goto end

:linux
echo build linux amd64
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-s -w" -o vnet  .\cmd\server\server.go
goto end

:end
echo done~