package natskv

import (
	"context"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sagikazarmark/crypt/backend"
)

type Client struct {
	conn   *nats.Conn
	js     nats.JetStreamContext
	bucket string
}

func New(machines []string) (*Client, error) {
	nc, err := nats.Connect(strings.Join(machines[:], ","))
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return &Client{conn: nc, js: js, bucket: "viper"}, nil

}

func (c *Client) Get(key string) ([]byte, error) {
	return c.GetWithContext(context.TODO(), key)
}

func (c *Client) GetWithContext(ctx context.Context, key string) ([]byte, error) {
	defer c.conn.Close()
	kv, err := c.js.KeyValue(c.bucket)
	if err != nil {
		return nil, err
	}

	v, err := kv.Get(key)
	if err != nil {
		return nil, err
	}

	return v.Value(), nil
}

func (c *Client) List(bucket string) (backend.KVPairs, error) {
	return c.ListWithContext(context.TODO(), bucket)
}

func (c *Client) ListWithContext(ctx context.Context, bucket string) (backend.KVPairs, error) {
	defer c.conn.Close()
	res := backend.KVPairs{}
	kv, err := c.js.KeyValue(c.bucket)
	if err != nil {
		return nil, err
	}

	ks, err := kv.Keys()
	if err != nil {
		return nil, err
	}

	for _, k := range ks {
		v, err := kv.Get(k)
		if err != nil {
			return nil, err
		}

		res = append(res, &backend.KVPair{
			Key:   v.Key(),
			Value: []byte(v.Value()),
		})
	}

	return res, nil

}

func (c *Client) Set(key string, value []byte) error {
	return c.SetWithContext(context.TODO(), key, value)
}

func (c *Client) SetWithContext(ctx context.Context, key string, value []byte) error {
	defer c.conn.Close()
	kv, err := c.js.KeyValue(c.bucket)
	if err != nil {
		return err
	}

	_, err = kv.Put(key, value)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Watch(key string, stop chan bool) <-chan *backend.Response {
	return c.WatchWithContext(context.TODO(), key, stop)
}

func (c *Client) WatchWithContext(ctx context.Context, key string, stop chan bool) <-chan *backend.Response {
	defer c.conn.Close()
	ch := make(chan *backend.Response, 0)

	kv, err := c.js.KeyValue(c.bucket)
	if err != nil {
		ch <- &backend.Response{nil, err}
	}

	watch, err := kv.Watch(key)
	if err != nil {
		ch <- &backend.Response{nil, err}
	}

	go func() {
		go func() {
			defer c.conn.Close()
			<-stop
			watch.Stop()
		}()

		for {
			k := <-watch.Updates()
			ch <- &backend.Response{
				Value: k.Value(),
			}
			time.Sleep(time.Second * 5)
		}

	}()

	return ch
}
