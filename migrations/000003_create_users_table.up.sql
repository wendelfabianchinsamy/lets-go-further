create table if not exists users (
    id bigserial primary key,
    email text not null,
    password text not null
);