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
-- +goose StatementBegin
SELECT create_distributed_table('dialogs', 'from_user_id');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
drop table if exists dialog;
-- +goose StatementEnd