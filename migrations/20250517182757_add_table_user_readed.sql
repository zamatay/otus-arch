-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_reade (
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    post_id INT REFERENCES posts(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    primary key (user_id, post_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists user_reade
-- +goose StatementEnd
