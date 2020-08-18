GOCMD=go
GOBUILD=$(GOCMD) build
BASEPATH := $(shell pwd)
BUILDDIR=$(BASEPATH)/dist

KOBE_SRC=$(BASEPATH)/cmd
KOBE_SERVER_NAME=kobe-server
KOBE_INVENTORY_NAME=kobe-inventory
KOBE_CLIENT_NAME=kobe

BIN_DIR=usr/local/bin
CONFIG_DIR=etc/kobe
BASE_DIR=var/kobe
LIB_DIR=$(BASE_DIR)/lib


build_server_linux:
	GOOS=linux  GOARCH=$(GOARCH) $(GOBUILD) -o $(BUILDDIR)/$(BIN_DIR)/$(KOBE_SERVER_NAME) $(KOBE_SRC)/server/*.go
	GOOS=linux  GOARCH=$(GOARCH) $(GOBUILD) -o $(BUILDDIR)/$(BIN_DIR)/$(KOBE_INVENTORY_NAME) $(KOBE_SRC)/inventory/*.go
	mkdir -p $(BUILDDIR)/$(LIB_DIR) && cp -r     $(BASEPATH)/ansible $(BUILDDIR)/$(LIB_DIR)
	mkdir -p $(BUILDDIR)/$(CONFIG_DIR) && cp -r  $(BASEPATH)/conf/* $(BUILDDIR)/$(CONFIG_DIR)

clean:
	rm -fr $(BUILDDIR)

docker:
	@echo "build docker images"
	docker build -t kubeoperator/kobe:master --build-arg GOARCH=$(GOARCH) .
