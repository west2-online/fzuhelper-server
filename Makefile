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

# 服务名
SERVICES := api user classroom
service = $(word 1, $@)

# mock gen
MOCKS := user_mock
mock = $(word 1, $@)

PERFIX = "[Makefile]"

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  $(SERVICES)   : Build a specific service (e.g., make api). use BUILD_ONLY=1 to avoid auto bootstrap."
	@echo "  env-up        : Start the docker-compose environment."
	@echo "  env-down      : Stop the docker-compose environment."
	@echo "  mocks         : Generate mocks for interfaces."
	@echo "  clean         : Remove the 'output' directories and related binaries."
	@echo "  clean-all     : Stop docker-compose services if running and remove 'output' directories and docker data."
	@echo "  docker        : Build a Docker image named 'fzuhelper'."

# 生成 mock 工具
.PHONY: $(MOCKS)
$(MOCKS):
	@mkdir -p mocks
	mockgen -source=./idl/$(mock).go -destination=./mocks/$(mock).go -package=mocks

# 启动必要的环境，比如 etcd、mysql
.PHONY: env-up
env-up:
	@ docker compose -f ./docker/docker-compose.yml up -d

# 关闭必要的环境，但不清理 data（位于 docker/data 目录中）
.PHONY: env-down
env-down:
	@ cd ./docker && docker compose down

# 生成 Kitex 相关代码
.PHONY: kitex-gen-%
kitex-gen-%:
	kitex \
	-gen-path ./kitex_gen \
	-service "$*" \
	-module "$(MODULE)" \
	-type thrift \
	./idl/$*.thrift
	go mod tidy

# 使用 Hertz 相关代码
hertz-gen-api:
	hz update -idl ./idl/api.thrift

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
		tmux select-layout -t "fzuhelper-$(service)" even-horizontal; \
	fi
	@echo "$(PERFIX) Running $(service) service in tmux..."
	@tmux send-keys -t fzuhelper-$(service).0 'sh entrypoint.sh $(service)' C-m
	@tmux select-pane -t fzuhelper-$(service).0
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