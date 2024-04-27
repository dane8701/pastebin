package store

import (
	"context"
	"encoding/json"
	"strings"
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
	clics := []ClicByBin{}
	stats := Statistics{}
	// tout calculer et retourner un json de mets stats
	for _, bin := range bins {
		clic := ClicByBin{BinID: bin.ID, Clic: bin.Clic}
		clics = append(clics, clic)
	}
	lenBins := len(bins)
	stats = Statistics{BinNumber: int32(lenBins), ClicByBin: clics}

	return &stats, nil
}

func (e *redisDB) CreateBin(ctx context.Context, bin Bin) (*Bin, error) {
	keys := e.client.Keys(ctx, "bin:"+bin.Alias+":*").Val()
	if len(keys) != 0 {
		return nil, errors.New(bin.Alias + " already exists as an alias")
	}

	bin.ID = uuid.NewString()
	value, err := json.Marshal(bin)
	if err != nil {
		return nil, errors.Wrapf(err, "couldnt json marshal bin %s", bin.ID)
	}

	expiration := 30 * 24 * time.Hour
	err = e.client.SetEx(ctx, redisKeyFrom(bin), string(value), expiration).Err()
	if err != nil {
		return nil, errors.Wrapf(err, "couldnt create bin %s", bin.ID)
	}

	return &bin, nil
}

func redisKeyFrom(bin Bin) string {
	if strings.TrimSpace(bin.Alias) == "" {
		return "bin:" + bin.ID
	}

	return "bin:" + bin.Alias + ":" + bin.ID
}

func (e *redisDB) GetBinByAlias(ctx context.Context, alias string) (*Bin, error) {
	if alias == "" {
		return nil, errors.Errorf("there is no alias provided")
	}

	keys := e.client.Keys(ctx, "bin:"+alias+":*").Val()
	if len(keys) == 0 {
		return nil, errors.New(alias + "couldnt query for bin")
	}
	val, err := e.client.Get(ctx, keys[0]).Result()
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

func (e *redisDB) SetBinExpiration(ctx context.Context, BD string, expiration time.Duration) error {
	key := "bin:" + BD
	return e.client.Expire(ctx, key, expiration).Err()
}

func (e *redisDB) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	if email == "" {
		return nil, errors.Errorf("there is no email provided")
	}

	keys := e.client.Keys(ctx, "user:"+email+":*").Val()
	if len(keys) == 0 {
			return nil, errors.New(email + "couldnt query for user")
	}
	val, err := e.client.Get(ctx, keys[0]).Result()
	if err != nil {
			return nil, errors.Wrapf(err, "couldnt query for user %s", email)
	}

	user := User{}
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse user from string")
	}

	return &user, nil
}

func (e *redisDB) CreateUser(ctx context.Context, user User) (*User, error) {
	userID := uuid.NewString()
	user.ID = userID

	//hâcher le mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.MotDePasse), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "failed to hash password")
	}
	user.MotDePasse = string(hashedPassword)

	userData, err := json.Marshal(user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal user data")
	}

	key := "user:"+user.Email+":"+userID
	err = e.client.Set(ctx, key, string(userData), 0).Err()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user in database")
	}

	return &user, nil
}

func (e *redisDB) GetAllUsers(ctx context.Context) ([]User, error) {
	keys, err := e.client.Keys(ctx, "user:*").Result()
	if err != nil {
			return nil, errors.Wrap(err, "couldnt query for users")
	}

	users := []User{}

	for _, id := range keys {
			val, err := e.client.Get(ctx, id).Result()
			if err != nil {
					return nil, errors.Wrapf(err, "couldnt query for user %s", id)
			}

			u := User{}
			err = json.Unmarshal([]byte(val), &u)
			if err != nil {
					return nil, errors.Wrap(err, "couldnt parsing user from string")
			}

			users = append(users, u)
	}

	return users, nil
}

func (e *redisDB) DropAllUsers(ctx context.Context) error {
	_, err := e.client.FlushDB(ctx).Result()
	if err != nil {
		return errors.Wrap(err, "failed to drop all users")
	}
	return nil
}