-- +goose Up
-- +goose StatementBegin
CREATE TABLE genders (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) UNIQUE NOT NULL
);

insert into genders(name)
values ('Мужской'), ('Женский');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login varchar(50) not null ,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    birthday DATE NULL,
    gender_id INT REFERENCES genders(id) ON DELETE SET NULL,
    interests TEXT[],
    city VARCHAR(100),
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE user_credentials (
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE user_credentials
DROP TABLE users;
DROP TABLE  genders
-- +goose StatementEnd
