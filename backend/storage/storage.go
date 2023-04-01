package storage

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/thisisheymde/URL-shortener/backend/types"
)

type Storage interface {
	InserttoDB(*types.Link) error
	GetfromDB(string) (string, error)
}

type redisStorage struct {
	client *redis.Client
}

var ctx = context.Background()

func StartRedis(addr string, password string) (*redisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	log.Println("database connected.")

	return &redisStorage{client: client}, nil
}

func (r *redisStorage) InserttoDB(l *types.Link) error {
	err := r.client.SetNX(ctx, l.ID, l.URL, time.Hour*24*7).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisStorage) GetfromDB(s string) (string, error) {
	val, err := r.client.Get(ctx, s).Result()
	if err != nil {
		return "", errors.New("key doesn't exist")
	}
	return val, nil
}
