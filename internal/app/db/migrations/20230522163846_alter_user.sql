-- +goose Up
-- +goose StatementBegin
CREATE TYPE public.role AS ENUM
    ('user', 'admin');

ALTER TABLE IF EXISTS public.users
    ADD COLUMN role role;

INSERT INTO public.users(
	first_name, last_name, email, password,  status, role)
	VALUES ('Admin', 'User', 'shubham.aal@gmail.com', '$2a$10$.fJ56Xvu4IABh2ZbsvdPHOhKjPxR0SZxx1onjF25Ss.W/51qvZopW', true, 'admin');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
