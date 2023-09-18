MODULE = github.com/west2-online/fzuhelper-server

.PHONY: target
target:
	sh build.sh

.PHONY: clean
clean:
	@find . -type d -name "output" -exec rm -rf {} + -print