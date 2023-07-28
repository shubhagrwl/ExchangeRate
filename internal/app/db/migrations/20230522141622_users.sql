-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.users
(
    id bigserial NOT NULL,
    first_name text NOT NULL,
    last_name text,
    email text NOT NULL,
    password text NOT NULL,
    country_code text,
    phone text,
    address text,
    status boolean NOT NULL DEFAULT 'false',
    created_at timestamp without time zone NOT NULL DEFAULT 'NOW()',
    PRIMARY KEY (id),
    CONSTRAINT users_email_uk UNIQUE (email, phone)
        INCLUDE(email, phone)
);
    
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
