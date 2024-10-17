package client

import (
	"errors"
	"fmt"

	"github.com/elastic/go-elasticsearch"
	"github.com/west2-online/fzuhelper-server/config"
)

func NewEsClient() (*elasticsearch.Client, error) {
	if config.Elasticsearch == nil {
		return nil, errors.New("elasticsearch config is nil")
	}
	esConn := fmt.Sprintf("http://%s", config.Elasticsearch.Addr)
	cfg := elasticsearch.Config{
		Addresses: []string{esConn},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("Get es clint failed,error: %v", err)
	}
	return client, nil
}
