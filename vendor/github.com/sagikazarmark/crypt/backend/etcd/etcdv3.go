package etcd

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sagikazarmark/crypt/backend"
	"go.etcd.io/etcd/api/v3/mvccpb"
	goetcdv3 "go.etcd.io/etcd/client/v3"
)

type ClientV3 struct {
	ctx     context.Context
	timeout time.Duration
	client  *goetcdv3.Client
	keysAPI goetcdv3.KV
}

var CliTimeoutKey = "ETCDV3_CLI_TIMEOUT_SECOND"
var oneclim = map[string]*ClientV3{}
var mu sync.Mutex

func NewV3(machines []string) (*ClientV3, error) {
	return NewFromV3Config(goetcdv3.Config{
		Endpoints: machines,
	})
}

func NewFromV3Config(config goetcdv3.Config) (*ClientV3, error) {
	machines := config.Endpoints
	sort.Strings(machines)
	key := strings.Join(machines, ",")
	if cli, ok := oneclim[key]; !ok {
		mu.Lock()
		defer mu.Unlock()
		if cli2, ok2 := oneclim[key]; !ok2 {
			newClient, err := goetcdv3.New(config)
			if err != nil {
				return nil, fmt.Errorf("creating new etcd client for crypt.backend.Client: %v", err)
			}
			keysAPI := goetcdv3.NewKV(newClient)
			timeout := time.Second * 3
			tm := os.Getenv(CliTimeoutKey)
			if tm != "" {
				tmi, err := strconv.Atoi(tm)
				if err != nil {
					return nil, errors.Wrapf(err, "getting timeout param from env [%s], parse to number failed", tm)
				}
				timeout = time.Duration(tmi) * time.Second
			}
			c := &ClientV3{client: newClient, keysAPI: keysAPI, ctx: context.Background(), timeout: timeout}
			oneclim[key] = c
			return c, nil
		} else {
			return cli2, nil
		}
	} else {
		return cli, nil
	}
}

func (c *ClientV3) Get(key string) ([]byte, error) {
	tctx, cancelFunc := context.WithTimeout(c.ctx, c.timeout)
	defer cancelFunc()
	res, err := c.keysAPI.Get(tctx, key)
	if err != nil {
		return nil, err
	}
	if res.Count != 1 {
		return nil, fmt.Errorf("getting from etcd with key [%s], res count %d not equal to 1", key, res.Count)
	}
	return res.Kvs[0].Value, nil
}

func (c *ClientV3) List(key string) (backend.KVPairs, error) {
	tctx, cancelFunc := context.WithTimeout(c.ctx, c.timeout)
	defer cancelFunc()
	res, err := c.keysAPI.Get(tctx, key, goetcdv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var bk backend.KVPairs
	for _, kv := range res.Kvs {
		bk = append(bk, &backend.KVPair{
			Key:   string(kv.Key),
			Value: kv.Value,
		})
	}
	return bk, nil
}

func (c *ClientV3) Set(key string, value []byte) error {
	tctx, cancelFunc := context.WithTimeout(c.ctx, c.timeout)
	defer cancelFunc()
	_, err := c.keysAPI.Put(tctx, key, string(value))
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientV3) Watch(key string, stop chan bool) <-chan *backend.Response {
	respChan := make(chan *backend.Response, 0)
	cctx, cancelFunc := context.WithCancel(c.ctx)
	go func() {
		<-stop
		cancelFunc()
	}()
	go func() {
		wch := c.client.Watch(cctx, key, goetcdv3.WithPrevKV())
		for {
			select {
			case we := <-wch:
				for _, ev := range we.Events {
					switch ev.Type {
					case mvccpb.PUT:
						respChan <- &backend.Response{Value: ev.Kv.Value}
					case mvccpb.DELETE:
						//do nothing with delete event
						fmt.Println("find DETELE:", ev.PrevKv.Key, ev.PrevKv.Value)
					}
				}
			case <-cctx.Done():
				fmt.Println("stop watch")
				return
			}
		}
	}()
	return respChan
}
