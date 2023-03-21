-- create table users
CREATE TABLE IF NOT EXISTS users
(
    id             BIGSERIAL PRIMARY KEY,
    login          TEXT    NOT NULL,
    password       TEXT    NOT NULL,
    check_password BOOLEAN NOT NULL DEFAULT TRUE,
    is_admin       BOOLEAN NOT NULL DEFAULT FALSE,
    is_blocked     BOOLEAN NOT NULL DEFAULT FALSE
);

-- login unique constraint
ALTER TABLE users
    ADD CONSTRAINT users_login_ukey
        UNIQUE (login);

-- create table sessions
CREATE TABLE IF NOT EXISTS sessions
(
    token   TEXT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    expr    TIMESTAMP
);

-- foreign keys sessions -> users
ALTER TABLE sessions
    ADD CONSTRAINT sessions_user_id_fkey
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE;
