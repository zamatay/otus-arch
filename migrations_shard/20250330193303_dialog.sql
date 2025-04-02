-- +goose Up
-- +goose StatementBegin
CREATE TABLE dialogs (
    from_user_id INT,
    to_user_id INT,
    text varchar(4000) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
drop table if exists dialogs;
-- +goose StatementEnd