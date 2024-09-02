package service

import (
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/devworlds/eventlistener-redis-performance/internal/db/repository"
	"github.com/devworlds/eventlistener-redis-performance/internal/redis"
)

func CreateWallet(address string) error {
	err := repository.AddWallet(address)
	if err != nil {
		return err
	}

	rdb := redis.Connect()
	payload, _ := rdb.GetWallet()
	f := bloom.NewWithEstimates(1000000, 0.01)
	f.UnmarshalJSON([]byte(payload))
	f.Add([]byte(address))
	jsonFilter, _ := f.MarshalJSON()
	rdb.SetWallet(string(jsonFilter))
	rdb.PublishToChan()
	return nil
}
