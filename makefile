# Variables
GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME_CLIENT=client.exe
BINARY_NAME_SERVER=server.exe
CLIENT_DIR=botnet/client
SERVER_DIR=botnet/server

# Targets
all: client server

client:
	cd $(CLIENT_DIR)/ && $(GOBUILD) -o ../../$(BINARY_NAME_CLIENT) main.go

server:
	cd $(SERVER_DIR)/ && $(GOBUILD) -o ../../$(BINARY_NAME_SERVER) main.go