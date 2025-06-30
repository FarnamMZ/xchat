create table users (
    username varchar(50) unique,
    password varchar(255) not null,
    email varchar(100) not null unique,
    role int
);
-- DROP TABLE users;