package post

import (
	"context"
	"log/slog"
	"net/http"

	srvApi "github.com/zamatay/otus/arch/lesson-1/internal/api"
	"github.com/zamatay/otus/arch/lesson-1/internal/domain"
)

type FeedPosted interface {
	GetFeedCache(context.Context, int) ([]*domain.Post, bool)
	SetFeedCache(context.Context, int, []*domain.Post) error
}

const defaultLimit = 100

func (u *Post) Feed(writer http.ResponseWriter, request *http.Request) {
	ctx, done := srvApi.GetContext(request.Context())
	defer done()

	offset, _ := srvApi.GetByName(request, "offset")
	userId, _ := srvApi.GetByName(request, "user_id")
	limit, err := srvApi.GetByName(request, "limit")
	if err != nil {
		limit = defaultLimit
	}

	u.cache.GetFeedCache(ctx, userId)

	if posts, IsOk := u.cache.GetFeedCache(ctx, userId); IsOk && len(posts) > 0 {
		srvApi.SetOk(writer, posts)
		return
	}

	posts, err := u.service.FeedPost(ctx, offset, limit, userId)
	if err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}

	err = u.cache.SetFeedCache(ctx, userId, posts)
	if err != nil {
		slog.Error("Не удалось установить кэш", "error", err)
	}

	srvApi.SetOk(writer, posts)
}
