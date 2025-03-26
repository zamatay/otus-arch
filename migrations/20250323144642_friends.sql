-- +goose Up
-- +goose StatementBegin
CREATE TABLE friends (
                         from_user_id INT REFERENCES users(id) ON DELETE CASCADE,
                         to_user_id INT REFERENCES users(id) ON DELETE CASCADE,
                         created_at TIMESTAMP DEFAULT NOW(),
                         primary key (from_user_id, to_user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table friends;
-- +goose StatementEnd
