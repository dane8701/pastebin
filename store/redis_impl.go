package store

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type redisDB struct {
	client *redis.Client
}

func NewRedisDB(ctx context.Context, address string) (Store, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: address,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, errors.Wrap(err, "couldnt ping redis")
	}

	return &redisDB{
		client: rdb,
	}, nil
}

func (e *redisDB) GetAllBins(ctx context.Context) ([]Bin, error) {
	keys, err := e.client.Keys(ctx, "bin:*").Result()
	if err != nil {
		return nil, errors.Wrap(err, "couldnt query for bins")
	}

	bins := []Bin{}

	for _, id := range keys {
		val, err := e.client.Get(ctx, id).Result()
		if err != nil {
			return nil, errors.Wrapf(err, "couldnt query for bin %s", id)
		}

		t := Bin{}
		err = json.Unmarshal([]byte(val), &t)
		if err != nil {
			return nil, errors.Wrap(err, "couldnt parsing bins from string")
		}

		bins = append(bins, t)
	}

	return bins, nil
}

func (e *redisDB) GetStats(ctx context.Context) (*Statistics, error) {
	keys, err := e.client.Keys(ctx, "bin:*").Result()
	if err != nil {
		return nil, errors.Wrap(err, "couldnt query for bins")
	}

	bins := []Bin{}

	for _, id := range keys {
		val, err := e.client.Get(ctx, id).Result()
		if err != nil {
			return nil, errors.Wrapf(err, "couldnt query for bin %s", id)
		}

		t := Bin{}
		err = json.Unmarshal([]byte(val), &t)
		if err != nil {
			return nil, errors.Wrap(err, "couldnt parsing bins from string")
		}

		bins = append(bins, t)
	}
	stats := Statistics {}
	// tout calculer et retourner un json de mets stats

	return &stats, nil
}

func (e *redisDB) CreateBin(ctx context.Context, bin Bin) (*Bin, error) {
	bin.ID = uuid.NewString()

	value, err := json.Marshal(bin)
	if err != nil {
		return nil, errors.Wrapf(err, "couldnt json marshal bin %s", bin.ID)
	}

	err = e.client.Set(ctx, "bin:"+bin.ID, string(value), 0).Err()
	if err != nil {
		return nil, errors.Wrapf(err, "couldnt create bin %s", bin.ID)
	}

	return &bin, nil
}
func (e *redisDB) GetBinByAlias(ctx context.Context, alias string) (*Bin, error) {
	if alias == "" {
		return nil, errors.Errorf("there is no alias provided")
	}

	val, err := e.client.Get(ctx, "bin:"+alias).Result()
	if err != nil {
		return nil, errors.Wrapf(err, "couldnt query for bin %s", alias)
	}

	t := Bin{}
	err = json.Unmarshal([]byte(val), &t)
	if err != nil {
		return nil, errors.Wrap(err, "couldnt parsing bin from string")
	}

	//update alias for count the clic number
	t.Clic = t.Clic + 1
	value, err := json.Marshal(t)
	if err != nil {
		return nil, errors.Wrapf(err, "couldnt json marshal bin %s", t.ID)
	}

	err = e.client.Set(ctx, "bin:"+t.ID, string(value), 0).Err()
	if err != nil {
		return nil, errors.Wrapf(err, "couldnt update bin to count the clics %s", t.ID)
	}

	return &t, nil
}

func (e *redisDB) UpdateBin(ctx context.Context, bin Bin) (*Bin, error) {
	value, err := json.Marshal(bin)
	if err != nil {
		return nil, errors.Wrapf(err, "couldnt json marshal bin %s", bin.ID)
	}

	err = e.client.Set(ctx, "bin:"+bin.ID, string(value), 0).Err()
	if err != nil {
		return nil, errors.Wrapf(err, "couldnt update bin %s", bin.ID)
	}

	return &bin, nil
}

func (e *redisDB) DeleteBinByID(ctx context.Context, id string) (*Bin, error) {
	val, err := e.client.Get(ctx, "bin:"+id).Result()
	if err != nil {
		return nil, errors.Wrapf(err, "couldnt query for bin %s", id)
	}

	t := Bin{}
	err = json.Unmarshal([]byte(val), &t)
	if err != nil {
		return nil, errors.Wrap(err, "couldnt parsing bin from string")
	}

	err = e.client.Del(ctx, "bin:"+t.ID).Err()
	if err != nil {
		return nil, errors.Wrapf(err, "couldnt delete bin %s", t.ID)
	}

	return &t, nil
}
