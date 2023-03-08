CREATE TABLE IF NOT EXISTS user_groups
(
    id       SERIAL PRIMARY KEY,
    name     TEXT    NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    login         TEXT NOT NULL,
    password      TEXT NOT NULL,
    user_group_id bigint REFERENCES user_groups (id)
);

CREATE TABLE IF NOT EXISTS sessions
(
    token   TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expr    TIMESTAMP
);

alter table users
    add constraint users_user_groups_id_fk
        foreign key (user_group_id) references user_groups (id);

alter table sessions
    add foreign key (user_id) references users (id);
