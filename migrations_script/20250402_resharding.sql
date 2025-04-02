-- +goose NO TRANSACTION
-- +goose Up
-- +goose StatementBegin
select  citus_rebalance_start();
-- +goose StatementEnd