package ai

import (
	"github.com/west2-online/fzuhelper-server/config"

	"go.baoshuo.dev/llmfunc"
)

func NewFunction[T any, R any](
	handler llmfunc.OutputHandler[T, R],
	opts ...llmfunc.Option[*llmfunc.FunctionConfig],
) *llmfunc.Function[T, R] {
	client := llmfunc.NewClient(config.AI.Key, config.AI.Endpoint)
	return llmfunc.NewFunction[T, R](client, handler, opts...)
}
