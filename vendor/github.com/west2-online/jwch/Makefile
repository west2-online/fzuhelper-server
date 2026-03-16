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
