package client

import (
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/client"
)

func NewHzClient() (*client.Client, error) {
	cli, err := client.NewClient()
	if err != nil {
		return nil, fmt.Errorf("client.NewHzClient failed to NewClient: %w", err)
	}
	return cli, nil
}
