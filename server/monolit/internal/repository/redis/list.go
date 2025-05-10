package redis

import (
	"context"
)

func (c *Cache) GetLen(ctx context.Context, key string) int64 {
	cnt, _ := c.RDB.LLen(ctx, key).Result()
	return cnt
}
