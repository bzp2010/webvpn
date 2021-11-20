package storage

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/bzp2010/webvpn/internal/utils"
)

type EtcdStorage struct {
	client *clientv3.Client
}

func (s EtcdStorage) Get(ctx context.Context, key string) (string, error) {
	resp, err := s.client.Get(ctx, key)
	if err != nil {
		utils.Log().Errorf("etcd get failed: %s", err)
		return "", fmt.Errorf("etcd get failed: %s", err)
	}

	if resp.Count == 0 {
		utils.Log().Warningf("key: %s is not found", key)
		return "", fmt.Errorf("key: %s is not found", key)
	}

	return string(resp.Kvs[0].Value), nil
}

func (s EtcdStorage) List(ctx context.Context, key string) ([]Keypair, error) {
	resp, err := s.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		utils.Log().Errorf("etcd get failed: %s", err)
		return nil, fmt.Errorf("etcd get failed: %s", err)
	}

	var ret []Keypair
	for i := range resp.Kvs {
		key := string(resp.Kvs[i].Key)
		value := string(resp.Kvs[i].Value)

		data := Keypair{
			Key:   key,
			Value: value,
		}
		ret = append(ret, data)
	}

	return ret, nil
}

func (s EtcdStorage) Create(ctx context.Context, key, val string) error {
	_, err := s.client.Put(ctx, key, val)
	if err != nil {
		utils.Log().Errorf("etcd put failed: %s", err)
		return fmt.Errorf("etcd put failed: %s", err)
	}
	return nil
}

func (s EtcdStorage) Update(ctx context.Context, key, val string) error {
	_, err := s.client.Put(ctx, key, val)
	if err != nil {
		utils.Log().Errorf("etcd put failed: %s", err)
		return fmt.Errorf("etcd put failed: %s", err)
	}
	return nil
}

func (s EtcdStorage) BatchDelete(ctx context.Context, keys []string) error {
	for i := range keys {
		resp, err := s.client.Delete(ctx, keys[i])
		if err != nil {
			utils.Log().Errorf("delete etcd key[%s] failed: %s", keys[i], err)
			return fmt.Errorf("delete etcd key[%s] failed: %s", keys[i], err)
		}

		if resp.Deleted == 0 {
			utils.Log().Warningf("key: %s is not found", keys[i])
			return fmt.Errorf("key: %s is not found", keys[i])
		}
	}
	return nil
}

func (s EtcdStorage) Watch(ctx context.Context, key string) <-chan WatchEvent {
	panic("Not yet implement")
}
