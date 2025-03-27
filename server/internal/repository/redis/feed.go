package redis

import (
	"context"
	"encoding/json"
	"strconv"

	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

func (c *Cache) GetFeedKey(userId int) string {
	return "feed:" + strconv.Itoa(userId)
}

func (c *Cache) GetFeedCache(ctx context.Context, userId int) ([]*domain.Post, bool) {
	key := c.GetFeedKey(userId)
	if cnt, _ := c.RDB.LLen(ctx, key).Result(); cnt > 0 {
		result, err := c.RDB.LRange(ctx, key, 0, cnt).Result()
		if err != nil {
			return nil, false
		}
		posts := make([]*domain.Post, 0, len(result))
		for _, r := range result {
			p := domain.Post{}
			err := json.Unmarshal([]byte(r), &p)
			if err != nil {
				return nil, false
			}
			posts = append(posts, &p)
		}
		return posts, true
	}
	return nil, false
}

func (c *Cache) SetFeedCache(ctx context.Context, userId int, posts []*domain.Post) error {
	return c.InsertPosts(ctx, posts)
}

func (c *Cache) DelLeftFeedCache(ctx context.Context, userId int, count int) error {
	key := c.GetFeedKey(userId)
	err := c.RDB.LPopCount(ctx, key, count).Err()
	return err
}

func (c *Cache) InsertPosts(ctx context.Context, posts []*domain.Post) error {
	p := c.RDB.Pipeline()
	for _, post := range posts {
		key := c.GetFeedKey(post.UserID)
		marshal, err := json.Marshal(post)
		if err != nil {
			return err
		}
		if err := p.LPush(ctx, key, marshal).Err(); err != nil {
			return err
		}
	}
	_, err := p.Exec(ctx)
	return err
}
