DIR = $(shell pwd)
CMD = $(DIR)/cmd
CONFIG_PATH = $(DIR)/config
IDL_PATH = $(DIR)/idl
OUTPUT_PATH = $(DIR)/output
SERVICES := template empty_room user api launch_screen
service = $(word 1, $@)

# mock gen
MOCKS := user_mock
mock = $(word 1, $@)

# hz&kitex
RPC = $(DIR)/cmd
API_PATH= $(DIR)/cmd/api
MODULE=github.com/west2-online/fzuhelper-server
KITEX_GEN_PATH=$(DIR)/kitex_gen

PERFIX = "[Makefile]"

.PHONY: env-up
env-up:
	@ docker compose -f ./docker/docker-compose.yml up -d

.PHONY: env-down
env-down:
	@ cd ./docker && docker compose down

# build specific target
.PHONY: $(SERVICES)
$(SERVICES):
	mkdir -p output
	cd $(CMD)/$(service) && sh build.sh
	cd $(CMD)/$(service)/output && cp -r . $(OUTPUT_PATH)/$(service)
	@echo "$(PERFIX) Build $(service) target completed"
ifndef autorun
	@echo "$(PERFIX) Automatic run server"
	sh entrypoint.sh $(service)
endif

# mock
.PHONY: $(MOCKS)
$(MOCKS):
	@mkdir -p mocks
	mockgen -source=./idl/$(mock).go -destination=./mocks/$(mock).go -package=mocks

# clean targets
.PHONY: clean
clean:
	@find . -type d -name "output" -exec rm -rf {} + -print

# build all
.PHONY: build-all
build-all:
	@for var in $(SERVICES); do \
		echo "$(PERFIX) building $$var service"; \
		make $$var autorun=1; \
	done

# use docker instead to run projects
.PHONY: docker
docker:
	docker build -t fzuhelper .

KSERVICES := user launch_screen
.PHONY: kgen
kgen:
	@for kservice in $(KSERVICES); do \
		kitex -module ${MODULE} ${IDL_PATH}/$$kservice.thrift; \
    	cd ${RPC};cd $$kservice;kitex -module ${MODULE} -service $$kservice -use ${KITEX_GEN_PATH} ${IDL_PATH}/$$kservice.thrift; \
    	cd ../../; \
    done \


.PHONY: hzgen
hzgen:
	cd ${API_PATH}; \
	hz update -idl ${IDL_PATH}/api.thrift; \
	swag init; \

