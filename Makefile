DIR = $(shell pwd)
CMD = $(DIR)/cmd
CONFIG_PATH = $(DIR)/config
IDL_PATH = $(DIR)/idl
OUTPUT_PATH = $(DIR)/output

SERVICES := template empty_room
service = $(word 1, $@)

# mock gen
MOCKS := user_mock
mock = $(word 1, $@)

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