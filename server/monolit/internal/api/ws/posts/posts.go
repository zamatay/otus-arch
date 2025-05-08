package posts

import (
	"githib.com/zamatay/otus/arch/lesson-1/internal/api/post"
)

type Posts struct {
	srv post.PostServiced
}

func NewPosts(srv post.PostServiced) *Posts {
	return &Posts{srv: srv}
}
