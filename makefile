# Variables
GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME_CLIENT_WINDOWS=client.exe
BINARY_NAME_CLIENT_LINUX=client
BINARY_NAME_SERVER_WINDOWS=server.exe
BINARY_NAME_SERVER_LINUX=server
CLIENT_DIR=botnet/client
SERVER_DIR=botnet/server

build: client_windows client_linux server_windows server_linux

client_windows:
	cd $(CLIENT_DIR)/ && cmd /C "set GOOS=windows&& set GOARCH=amd64&& $(GOBUILD) -o ../../$(BINARY_NAME_CLIENT_WINDOWS) main.go"

client_linux:
	cd $(CLIENT_DIR)/ && cmd /C "set GOOS=linux&& set GOARCH=amd64&& $(GOBUILD) -o ../../$(BINARY_NAME_CLIENT_LINUX) main.go"

server_windows:
	cd $(SERVER_DIR)/ && cmd /C "set GOOS=windows&& set GOARCH=amd64&& $(GOBUILD) -o ../../$(BINARY_NAME_SERVER_WINDOWS) main.go"

server_linux:
	cd $(SERVER_DIR)/ && cmd /C "set GOOS=linux&& set GOARCH=amd64&& $(GOBUILD) -o ../../$(BINARY_NAME_SERVER_LINUX) main.go"

vm_space:
	@echo [+] starting up virtual machines
	VBoxManage startvm "windows1VM"
	VBoxManage startvm "windows2VM"
	VBoxManage startvm "linux1VM"
	VBoxManage startvm "linux2VM"

deploy: build
	@echo [+] deploying code to virtual machines
	cd vmconfig && vagrant plugin install vagrant-scp
	
	cd vmconfig && vagrant scp ../client.exe W1VM:/C:/Users/vagrant/desktop/
	cd vmconfig && vagrant scp ../.env W1VM:/C:/Users/vagrant/desktop/

	cd vmconfig && vagrant scp ../client.exe W2VM:/C:/Users/vagrant/desktop/
	cd vmconfig && vagrant scp ../.env W2VM:/C:/Users/vagrant/desktop/

	cd vmconfig && vagrant scp ../client linux1VM:/home/vagrant
	cd vmconfig && vagrant scp ../.env linux1VM:/home/vagrant

	cd vmconfig && vagrant scp ../client linux2VM:/home/vagrant
	cd vmconfig && vagrant scp ../.env linux2VM:/home/vagrant
	@echo [+] done...

help:
	@echo [+] deploy :: build and deploy client and .env to vmspace
	@echo [+] vm_space :: start all the VM
	@echo [+] build :: compile the server and client to executables