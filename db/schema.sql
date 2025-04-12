-- with postgres user
CREATE DATABASE vto;
CREATE USER vto WITH PASSWORD 'vto';
GRANT ALL PRIVILEGES ON DATABASE vto TO vto;

\c vto

GRANT ALL ON SCHEMA public TO vto;

-- with vto user

CREATE TABLE shop (
    id BIGSERIAL PRIMARY KEY,
    access_token TEXT NOT NULL,
    shop_domain TEXT NOT NULL UNIQUE
);