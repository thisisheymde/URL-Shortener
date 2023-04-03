package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

// var clientRL = redis.NewClient(&redis.Options{
// 	Addr:     os.Getenv("REDISCACHE_HOST") + ":" + os.Getenv("REDISCACHE_PORT"),
// 	Password: os.Getenv("REDISCACHE_PASSWORD"),
// })

var clientRL = redis.NewClient(&redis.Options{
	Addr:     "containers-us-west-163.railway.app:6693",
	Password: "RZACxSmmqZVvhXUYVgfu",
})

var ctx = context.Background()

func RateLimiting(w http.ResponseWriter, r *http.Request) error {
	exists, _ := clientRL.Exists(ctx, r.RemoteAddr).Result()

	if exists == 0 {
		err := clientRL.Set(ctx, r.RemoteAddr, 1, time.Minute).Err()
		if err != nil {
			return err
		}
		return nil
	}

	count, err := clientRL.Get(ctx, r.RemoteAddr).Int64()
	if err != nil {
		return err
	}

	if count > 25 {
		return errors.New("rate exceeded")
	}

	err = clientRL.Set(ctx, r.RemoteAddr, count+1, redis.KeepTTL).Err()
	if err != nil {
		return err
	}

	return nil
}
