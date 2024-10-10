# 默认输出帮助信息
.DEFAULT_GOAL := help
# 检查 tmux 是否存在
TMUX_EXISTS := $(shell command -v tmux)
# 远程仓库
REMOTE_REPOSITORY = registry.cn-hangzhou.aliyuncs.com/west2-online/fzuhelper-server
# 项目 MODULE 名
MODULE = github.com/west2-online/fzuhelper-server
# 当前架构
ARCH := $(shell uname -m)
# 目录相关
DIR = $(shell pwd)
CMD = $(DIR)/cmd
CONFIG_PATH = $(DIR)/config
IDL_PATH = $(DIR)/idl
OUTPUT_PATH = $(DIR)/output
API_PATH= $(DIR)/cmd/api

# 服务名
SERVICES := api user classroom launch_screen
service = $(word 1, $@)

PERFIX = "[Makefile]"

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  {service name}    : Build a specific service (e.g., make api). use BUILD_ONLY=1 to avoid auto bootstrap."
	@echo "                      Available service list: [${SERVICES}]"
	@echo "  env-up            : Start the docker-compose environment."
	@echo "  env-down          : Stop the docker-compose environment."
	@echo "  kitex-gen-%       : Generate Kitex service code for a specific service (e.g., make kitex-gen-user)."
	@echo "  kitex-update-%    : Update Kitex generated code for a specific service (e.g., make kitex-update-user)."
	@echo "  hertz-gen-api     : Generate Hertz scaffold based on the API IDL."
	@echo "  test              : Run unit tests for the project."
	@echo "  clean             : Remove the 'output' directories and related binaries."
	@echo "  clean-all         : Stop docker-compose services if running and remove 'output' directories and docker data."
	@echo "  push-%            : Push a specific service to the remote repository (e.g., make push-api)."
	@echo "  fmt               : Format the codebase using gofumpt."
	@echo "  import            : Optimize import order and structure."
	@echo "  vet               : Check for possible errors with go vet."
	@echo "  lint              : Run golangci-lint on the codebase."
	@echo "  verify            : Format, optimize imports, and run linters and vet on the codebase."
	@echo "  license           : Check and add license to go file and shell script."

## --------------------------------------
## 构建与调试
## --------------------------------------

# 启动必要的环境，比如 etcd、mysql
.PHONY: env-up
env-up:
	@ docker compose -f ./docker/docker-compose.yml up -d

# 关闭必要的环境，但不清理 data（位于 docker/data 目录中）
.PHONY: env-down
env-down:
	@ cd ./docker && docker compose down

# 生成基于 Kitex 的业务代码，在新建业务时使用
# TODO: 这么写是因为 kitex 这个 cli 太难用了，计划修改成 cwgo 的
.PHONY: kitex-gen-%
kitex-gen-%:
	mkdir -p $(CMD)/$* && cd $(CMD)/$* && \
	kitex \
	-gen-path ../../kitex_gen \
	-service "$*" \
	-module "$(MODULE)" \
	-type thrift \
	$(DIR)/idl/$*.thrift
	go mod tidy

# 更新 kitex_gen 中的对应模块，不会影响 cmd 中的业务
.PHONY: kitex-update-%
kitex-update-%:
	kitex -module "${MODULE}" idl/$*.thrift

# 生成基于 Hertz 的脚手架
# TODO: 这个和 Kitex 的区别在于这个 update 实际上做了 gen 的工作，就直接这么用了
.PHONY: hertz-gen-api
hertz-gen-api:
	cd ${API_PATH}; \
    hz update -idl ${IDL_PATH}/api.thrift; \
	#hz model -idl ./idl/api.thrift  --model_dir ./cmd/api/biz/model && \
#	hz update -idl ./idl/api.thrift \
#	--out_dir ./cmd/api  \
#	--use ${MODULE}/cmd/api/biz/model  \
#	--handler_dir ./biz/handler \
#	--router_dir ./biz/router && \
#	使用其会导致handler直接丢失，router的import错误
	cd $(DIR) && \
	swag init --dir ./cmd/api --output ./docs/swagger --outputTypes "yaml"
# 单元测试
.PHONY: test
test:
	go test -v -gcflags="all=-l -N" -coverprofile=coverage.txt -parallel=16 -p=16 -covermode=atomic -race -coverpkg=./... \
		`go list ./... | grep -E -v "kitex_gen|.github|idl|docs|config|deploy"`

# 构建指定对象，构建后在没有给 BUILD_ONLY 参的情况下会自动运行，需要熟悉 tmux 环境
# 用于本地调试
.PHONY: $(SERVICES)
$(SERVICES):
	@if [ -z "$(TMUX_EXISTS)" ]; then \
		echo "$(PERFIX) tmux is not installed. Please install tmux first."; \
		exit 1; \
	fi
	@if [ -z "$$TMUX" ]; then \
		echo "$(PERFIX) you are not in tmux, press ENTER to start tmux environment."; \
		read -r; \
		if tmux has-session -t fzuhelp 2>/dev/null; then \
			echo "$(PERFIX) Tmux session 'fzuhelp' already exists. Attaching to session and running command."; \
			tmux attach-session -t fzuhelp; \
			tmux send-keys -t fzuhelp "make $(service)" C-m; \
		else \
			echo "$(PERFIX) No tmux session found. Creating a new session."; \
			tmux new-session -s fzuhelp "make $(service); $$SHELL"; \
		fi; \
	else \
		echo "$(PERFIX) Build $(service) target..."; \
		mkdir -p output; \
		cd $(CMD)/$(service) && sh build.sh; \
		cd $(CMD)/$(service)/output && cp -r . $(OUTPUT_PATH)/$(service); \
		echo "$(PERFIX) Build $(service) target completed"; \
	fi
ifndef BUILD_ONLY
	@echo "$(PERFIX) Automatic run server"
	@if tmux list-windows -F '#{window_name}' | grep -q "^fzuhelper-$(service)$$"; then \
		echo "$(PERFIX) Window 'fzuhelper-$(service)' already exists. Reusing the window."; \
		tmux select-window -t "fzuhelper-$(service)"; \
	else \
		echo "$(PERFIX) Window 'fzuhelper-$(service)' does not exist. Creating a new window."; \
		tmux new-window -n "fzuhelper-$(service)"; \
		tmux split-window -h ; \
		tmux select-layout -t "fzuhelper-$(service)" even-horizontal; \
	fi
	@echo "$(PERFIX) Running $(service) service in tmux..."
	@tmux send-keys -t fzuhelper-$(service).0 'sh ./hack/entrypoint.sh $(service)' C-m
	@tmux select-pane -t fzuhelper-$(service).1
endif

# 推送到镜像服务中，需要提前 docker login，否则会推送失败
# 不设置同时推送全部服务，这是一个非常危险的操作
.PHONY: push-%
push-%:
	@read -p "Confirm service name to push (type '$*' to confirm): " CONFIRM_SERVICE; \
	if [ "$$CONFIRM_SERVICE" != "$*" ]; then \
		echo "Confirmation failed. Expected '$*', but got '$$CONFIRM_SERVICE'."; \
		exit 1; \
	fi; \
	@if echo "$(SERVICES)" | grep -wq "$*"; then \
		if [ "$(ARCH)" = "x86_64" ] || [ "$(ARCH)" = "amd64" ] ; then \
			echo "Building and pushing $* for amd64 architecture..."; \
			docker build --build-arg SERVICE=$* -t $(REMOTE_REPOSITORY):$* -f docker/Dockerfile .; \
			docker push $(REMOTE_REPOSITORY):$*; \
		else \
			echo "Building and pushing $* using buildx for amd64 architecture..."; \
			docker buildx build --platform linux/amd64 --build-arg SERVICE=$* -t $(REMOTE_REPOSITORY):$* -f docker/Dockerfile --push .; \
		fi; \
	else \
		echo "Service '$*' is not a valid service. Available: [$(SERVICES)]"; \
		exit 1; \
	fi

## --------------------------------------
## 清理与校验
## --------------------------------------

# 清除所有的构建产物
.PHONY: clean
clean:
	@find . -type d -name "output" -exec rm -rf {} + -print

# 清除所有构建产物、compose 环境和它的数据
.PHONY: clean-all
clean-all: clean
	@echo "$(PREFIX) Checking if docker-compose services are running..."
	@docker-compose -f ./docker/docker-compose.yml ps -q | grep '.' && docker-compose -f ./docker/docker-compose.yml down || echo "$(PREFIX) No services are running."
	@echo "$(PREFIX) Removing docker data..."
	rm -rf ./docker/data

# 格式化代码，我们使用 gofumpt，是 fmt 的严格超集
.PHONY: fmt
fmt:
	gofumpt -l -w .

# 优化 import 顺序结构
.PHONY: import
import:
	goimports -w -local github.com/west2-online .

# 检查可能的错误
.PHONY: vet
vet:
	go vet ./...

# 代码格式校验
.PHONY: lint
lint:
	golangci-lint run --config=./.golangci.yml

# 一键修正规范并执行代码检查
.PHONY: verify
verify: vet fmt import lint

# 补齐 license
.PHONY: license
license:
	sh ./hack/add-license.sh
