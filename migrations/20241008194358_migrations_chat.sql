-- +goose Up
create table chat (
  id serial primary key,
  user_names text[], -- Array of text values
  -- user_from text not null,
  -- user_text text not null,
  send_time timestamp not null default now()
);



-- +goose Down
drop table chat;