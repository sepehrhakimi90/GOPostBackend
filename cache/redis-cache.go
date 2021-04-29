package cache

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis"

	"github.com/sepehrhakimi90/GOPostBackend/entity"
)

type redisCache struct {
	host    string
	db      int
	expires time.Duration
}

func (r *redisCache) Set(key string, value *entity.Post) {
	client := r.getClient()

	js, err := json.Marshal(value)
	if err != nil {
		log.Println("Cache", err)
		panic(err)
	}

	client.Set(key, js, r.expires)
}

func (r *redisCache) Get(key string) *entity.Post {
	client := r.getClient()

	val, err := client.Get(key).Result()
	if err != nil {
		log.Println("Cache", err)
		return nil
	}

	post := entity.Post{}
	err = json.Unmarshal([]byte(val), &post)
	if err != nil {
		log.Println("Cache", err)
		return nil
	}

	return &post
}

func NewRedisCache(host string, db int, exp time.Duration) PostCache {
	return &redisCache{
		host:    host,
		db:      db,
		expires: exp,
	}
}

func (r *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     r.host,
		Password: "",
		DB:       r.db,
	})
}
