-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.exchange
(
    id bigserial NOT NULL,
    exchange_date date NOT NULL,
    base text,
    usd decimal,
    eur decimal,
    gbp decimal,
    created_at timestamp without time zone NOT NULL DEFAULT 'NOW()',
    PRIMARY KEY (id),
    CONSTRAINT exchange_date_uk UNIQUE (exchange_date)
        INCLUDE(exchange_date)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS exchange;
-- +goose StatementEnd
