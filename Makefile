DIR = $(shell pwd)
CMD = $(DIR)/cmd
CONFIG_PATH = $(DIR)/config
IDL_PATH = $(DIR)/idl
OUTPUT_PATH = $(DIR)/output

SERVICES := api user classroom
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
	cd docker && docker build -t fzuhelper .

#允许传入特定服务进行构建，例如：make docker-build SERVICE=api
.PHONY: docker-build
docker-build:
	@if [ -z "$(SERVICE)" ]; then \
		for service in $(SERVICES); do \
			echo "Building Docker image for $$service..."; \
			docker build --build-arg SERVICE=$${service} -t $${service}_image -f docker/Dockerfile .; \
		done \
	else \
		echo "Building Docker image for $(SERVICE)..."; \
		docker build --build-arg SERVICE=$${SERVICE} -t $${SERVICE}_image -f docker/Dockerfile .; \
	fi


#启动所有服务
.PHONY: deploy
deploy:
	@ sh ./deploy/start-service-all.sh

#停止所有服务
.PHONY: stop
stop:
	for service in $(SERVICES); do \
  		docker rm -f $${service}; \
	done
