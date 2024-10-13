-- +goose Up
create table chat (
    id serial primary key,
    userName text not null,
    userFrom text not null,
    userText text not null,
    sendTime timestamp not null default now(),
);

-- +goose Down
drop table chat;