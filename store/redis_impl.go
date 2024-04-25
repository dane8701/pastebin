package store

import (
	"context"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
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

	expiration := 30 * 24 * time.Hour
	err = e.SetBinExpiration(ctx, bin.ID, expiration)
	if err != nil {
			return nil, errors.Wrapf(err, "couldnt set expiration for bin %s", bin.ID)
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

func (e *redisDB) SetBinExpiration(ctx context.Context, binID string, expiration time.Duration) error {
	key := "bin:" + binID
	return e.client.Expire(ctx, key, expiration).Err()
}

func (e *redisDB) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	if email == "" {
			return nil, errors.Errorf("there is no email provided")
	}

	val, err := e.client.Get(ctx, "user:"+email).Result()
	if err != nil {
			return nil, errors.Wrapf(err, "couldnt query for user %s", email)
	}

	var user User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
			return nil, errors.Wrap(err, "couldnt parsing user from string")
	}

	return &user, nil
}

func (e *redisDB) CreateUser(ctx context.Context, user User) (*User, error) {
	userID := uuid.NewString()
	user.ID = userID

	// Hacher le mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.MotDePasse), bcrypt.DefaultCost)
	if err != nil {
			return nil, errors.Wrap(err, "failed to hash password")
	}
	user.MotDePasse = string(hashedPassword)

	userData, err := json.Marshal(user)
	if err != nil {
			return nil, errors.Wrap(err, "failed to marshal user data")
	}

	key := "user:" + userID
	err = e.client.Set(ctx, key, string(userData), 0).Err()
	if err != nil {
			return nil, errors.Wrap(err, "failed to create user in database")
	}

	return &user, nil
}