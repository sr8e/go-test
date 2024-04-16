create table ir_user (
    id varchar not null primary key,
    display_name varchar not null,
    icon_url varchar not null,
    access_token varchar not null,
    refresh_token varchar not null,
    expire timestamptz not null,
    secret_hash varchar not null,
    secret_salt varchar not null
);