SHELL=/bin/bash -o pipefail

BINARY_NAME=cxa-ir
SERVICE_NAME=$(BINARY_NAME).service
SERVICE_TEMPLATE=systemd.service.template
SERVICE_PATH=/etc/systemd/system/$(SERVICE_NAME)
BINARY_DIR=$(shell pwd)/bin
BINARY_PATH=$(BINARY_DIR)/$(BINARY_NAME)
WORKING_DIR=$(shell pwd)

all: build

build:
	@mkdir -p bin
	go build -o ./bin/$(BINARY_NAME) .
	sudo setcap 'cap_net_bind_service=+ep' ./bin/$(BINARY_NAME)

build_test: build
	./bin/$(BINARY_NAME)

install: check_irsend
	sudo cp -vi LIRC-Remote/*.conf /etc/lirc/lircd.conf.d/
	sudo systemctl restart lircd.service

	@sed \
		-e 's|{{BINARY_PATH}}|$(BINARY_PATH)|' \
		-e 's|{{BINARY_NAME}}|$(BINARY_NAME)|' \
		-e 's|{{WORKING_DIR}}|$(WORKING_DIR)|' \
		-e 's|{{USER}}|$(USER)|' \
		-e 's|{{GROUP}}|$(USER)|' \
		$(SERVICE_TEMPLATE) | sudo tee $(SERVICE_PATH) > /dev/null
	sudo systemctl daemon-reload
	sudo systemctl enable $(SERVICE_NAME)
	sudo systemctl restart $(SERVICE_NAME)
	@sleep 1
	PAGER=cat systemctl status $(SERVICE_NAME)

check_irsend:
	@if ! command -v irsend &> /dev/null; then \
		echo "irsend could not be found. Please install lirc."; \
		exit 1; \
	fi
	@echo "irsend: OK"

clean:
	sudo systemctl stop $(SERVICE_NAME) 2>/dev/null || true
	sudo systemctl disable $(SERVICE_NAME) 2>/dev/null || true
	sudo rm -f $(SERVICE_PATH) 2>/dev/null || true
	@sudo systemctl daemon-reload
	@rm -rf bin

.PHONY: all build clean
