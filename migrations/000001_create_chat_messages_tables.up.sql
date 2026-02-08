create table chats (
    id serial primary key,
    title varchar(255) not null,
    created_at timestamp with time zone default now()
);

create table messages (
    id serial primary key,
    chat_id integer not null references chats(id) on delete cascade,
    text text not null,
    created_at timestamp with time zone default now()
);
