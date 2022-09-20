package sqlc

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kanengo/goutil/pkg/utils"
	"gorm.io/gorm"
)

type (
	ExecCtxFn  func(ctx context.Context, conn *gorm.DB) error
	QueryCtxFn func(ctx context.Context, conn *gorm.DB, v interface{}) error
)

type CachedConn struct {
	db           *gorm.DB
	cache        *redis.Client
	cacheTimeout time.Duration
}

func NewConn(db *gorm.DB, cache *redis.Client) CachedConn {
	return CachedConn{
		db:           db,
		cache:        cache,
		cacheTimeout: 5 * time.Minute,
	}
}

func NewConnWithCacheTimeout(db *gorm.DB, cache *redis.Client, timeout time.Duration) CachedConn {
	return CachedConn{
		db:           db,
		cache:        cache,
		cacheTimeout: timeout,
	}
}

func (cc CachedConn) DelCacheCtx(ctx context.Context, keys ...string) error {
	if cc.cache == nil {
		return nil
	}
	return cc.cache.Del(ctx, keys...).Err()
}

func (cc CachedConn) ExecCtx(ctx context.Context, exec ExecCtxFn, keys ...string) error {
	err := exec(ctx, cc.db)
	if err != nil {
		return err
	}

	if err := cc.DelCacheCtx(ctx, keys...); err != nil {
		return err
	}

	return nil
}

func (cc CachedConn) FindOneCtx(ctx context.Context, v any, key string, query QueryCtxFn) error {
	if cc.cache != nil {
		cacheRes, err := cc.cache.Get(ctx, key).Result()
		if err != nil && err != redis.Nil {
			return err
		}
		if err == nil {
			err = json.Unmarshal(utils.StringToSliceBytesUnsafe(cacheRes), v)
			if err == nil {
				return nil
			}
		}
	}

	if err := query(ctx, cc.db, v); err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	if cc.cache != nil {
		cc.cache.Set(ctx, key, v, cc.cacheTimeout)
	}

	return nil
}
