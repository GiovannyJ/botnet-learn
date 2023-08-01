# Variables
GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME_CLIENT_WINDOWS=client.exe
BINARY_NAME_CLIENT_LINUX=client
BINARY_NAME_SERVER_WINDOWS=server.exe
BINARY_NAME_SERVER_LINUX=server
CLIENT_DIR=botnet/client
SERVER_DIR=botnet/server

# Targets
all: client_windows client_linux server_windows server_linux

client_windows:
	cd $(CLIENT_DIR)/ && cmd /C "set GOOS=windows&& set GOARCH=amd64&& $(GOBUILD) -o ../../$(BINARY_NAME_CLIENT_WINDOWS) main.go"

client_linux:
	cd $(CLIENT_DIR)/ && cmd /C "set GOOS=linux&& set GOARCH=amd64&& $(GOBUILD) -o ../../$(BINARY_NAME_CLIENT_LINUX) main.go"

server_windows:
	cd $(SERVER_DIR)/ && cmd /C "set GOOS=windows&& set GOARCH=amd64&& $(GOBUILD) -o ../../$(BINARY_NAME_SERVER_WINDOWS) main.go"

server_linux:
	cd $(SERVER_DIR)/ && cmd /C "set GOOS=linux&& set GOARCH=amd64&& $(GOBUILD) -o ../../$(BINARY_NAME_SERVER_LINUX) main.go"
