include ../common.mk

SERVICE_NAME = classroom

.PHONY: server
server:
	sh build.sh
	cd output && sh bootstrap.sh

.PHONY: gen
gen:
	kitex \
	-gen-path ../../kitex_gen \
	-service "$(SERVICE_NAME)" \
	-module "$(MODULE)" \
	-type thrift \
	../../idl/$(SERVICE_NAME).thrift
	go mod tidy

.PHONY: test
test:
	cd test && go test -v

.PHONY: update
update:
	cd ../.. && kitex -module "${MODULE}" idl/"${SERVICE_NAME}".thrift
