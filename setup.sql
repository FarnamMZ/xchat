create table users (
    id serial primary key,
    username varchar(50) unique,
    password varchar(255) not null,
    email varchar(100) not null unique,
    role int
);

create table refresh_tokens (
    jti varchar(36) primary key,
    user_id int not null,
    expires_at timestamp not null,
    created_at timestamp default current_timestamp,
    foreign key (user_id) references users(id) on delete cascade
);

-- DROP TABLE refresh_tokens;
-- DROP TABLE users;