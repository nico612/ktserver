package captcha

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

var _ Store = new(RedisStore)

type RedisStore struct {
	rds redis.UniversalClient
}

func NewRedisStore(rds redis.UniversalClient) *RedisStore {
	return &RedisStore{
		rds: rds,
	}
}

func (s *RedisStore) Set(ctx context.Context, key string, code CaptchaCode) error {
	sKey := storeKey(key)
	value, err := json.Marshal(code)
	if err != nil {
		return err
	}
	_, err = s.rds.Set(ctx, sKey, value, code.Expiration).Result()
	if err != nil {
		return err
	}
	return err
}

func (s *RedisStore) Get(ctx context.Context, key string) (CaptchaCode, error) {
	sKey := storeKey(key)
	val, err := s.rds.Get(ctx, sKey).Result()
	if err != nil {
		return CaptchaCode{}, err
	}

	var code CaptchaCode
	if err = json.Unmarshal([]byte(val), &code); err != nil {
		return CaptchaCode{}, err
	}

	return code, nil
}

func (s *RedisStore) Save(ctx context.Context, key string, value CaptchaCode) error {
	sKey := storeKey(key)
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = s.rds.Set(ctx, sKey, val, redis.KeepTTL).Result()
	if err != nil {
		return err
	}

	return nil
}

func (s *RedisStore) Delete(ctx context.Context, key string) error {
	sKey := storeKey(key)
	_, err := s.rds.Del(ctx, sKey).Result()
	if err != nil {
		return err
	}
	return nil
}
