package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"

	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

const keyTemplate = "posts:%d"
const keyTemplateIndex = "index"

func (c *Cache) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	pipeline := c.RDB.Pipeline()
	key := getKey(post.UserID)
	pipeline.HSet(ctx, key, getField(post))
	pipeline.HSet(ctx, keyTemplateIndex, map[string]any{post.ID: post.UserID})
	if _, err := pipeline.Exec(ctx); err != nil {
		return nil, err
	}
	//key := getKey(post.UserID)
	//result, err := c.RDB.HSet(ctx, key, getField(post)).Result()
	//if err != nil {
	//	return nil, err
	//}
	//fmt.Println(result)
	//result, err = c.RDB.HSet(ctx, keyTemplateIndex, map[string]any{post.ID: post.UserID}).Result()
	//if err != nil {
	//	return nil, err
	//}
	//fmt.Println(result)
	//str, err := c.RDB.HGet(ctx, key, post.ID).Result()
	//if err != nil {
	//	return nil, err
	//}
	//fmt.Println(str)
	//if _, err := pipeline.Exec(ctx); err != nil {
	//	return nil, err
	//}
	return post, nil
}

func getField(post *domain.Post) map[string]any {
	return map[string]any{
		post.ID: post.Text,
	}
}

func getKey(userId int) string {
	return fmt.Sprintf(keyTemplate, userId)
}
func (c *Cache) UpdatePost(ctx context.Context, post *domain.Post) (bool, error) {
	key := getKey(post.UserID)
	v := c.RDB.HExists(ctx, key, post.ID)
	if v.Err() != nil && !errors.Is(v.Err(), redis.Nil) {
		return false, v.Err()
	}
	if !v.Val() {
		return false, nil
	}
	if result := c.RDB.HMSet(ctx, key, getField(post)); result.Err() != nil {
		return true, nil
	}
	return false, nil
}

func (c *Cache) DeletePost(ctx context.Context, id string, userId int) (bool, error) {
	result := c.RDB.HDel(ctx, getKey(userId), id)
	if result.Err() != nil {
		return true, nil
	}
	return false, result.Err()
}

func (c *Cache) GetPost(ctx context.Context, id string) (*domain.Post, error) {
	result, err := c.RDB.HGet(ctx, keyTemplateIndex, id).Result()
	if err != nil {
		return nil, err
	}
	if result == "" {
		return nil, errors.New("post not found")
	}
	userId, err := strconv.Atoi(result)
	if err != nil {
		return nil, err
	}
	value, err := c.RDB.HMGet(ctx, getKey(userId), id).Result()
	if err != nil {
		return nil, err
	}
	if len(value) == 0 {
		return nil, errors.New("post not found")
	}
	res, ok := value[0].(string)
	if !ok {
		return nil, errors.New("post not found")
	}
	return &domain.Post{ID: id, UserID: userId, Text: res}, nil
}

func (c *Cache) FeedPost(ctx context.Context, offset int, limit int, userID int) ([]*domain.Post, error) {
	cursor := uint64(0)
	var err error
	result := make([]*domain.Post, 0, limit)
	var strResult []string
	countPage := offset / limit
	count := 0
	for {
		strResult, cursor, err = c.RDB.HScan(ctx, getKey(userID), cursor, "*", int64(limit)).Result()
		if err != nil {
			return nil, err
		}

		if countPage == count {
			for i := 0; i < len(strResult); i += 2 {
				result = append(result, &domain.Post{ID: strResult[i], UserID: userID, Text: strResult[i+1]})
			}
		}

		if cursor == 0 {
			break
		}
		count++
	}
	return result, nil
}
