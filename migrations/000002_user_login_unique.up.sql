alter table users
    add constraint users_login_ukey
        unique (login);