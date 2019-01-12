create database siteforeg;
use siteforeg;

create table locations (
    `id` int auto_increment primary key not null, 
    `name` varchar(255) not null default "", 
    `address` varchar(255) not null default ""
    );

create table roles (
    `id` int primary key not null,
    `name` varchar(100) not null default "",
    `lvl` int not null default 1
);

INSERT INTO roles (id,name,lvl) VALUES (1,"root", 4);
-- есть идея что остальные роли не нужны
-- INSERT INTO roles (id,name,lvl) VALUES (2,"organizer", 3);
-- INSERT INTO roles (id,name,lvl) VALUES (3,"lector", 2);
-- INSERT INTO roles (id,name,lvl) VALUES (4,"listener", 1);


create table users (
    `id` int auto_increment primary key not null, 
    `password` varchar(30) not null default "", 
    `email` varchar(100) not null default "",
    `is_email_confirmed` bool not null default 0,
    `confirm_secret` varchar(300) not null default "",
    `fio` varchar(250) not null default "",
    `role_id` int not null default 4 
    );

-- INSERT INTO users (id, password, email, is_email_confirmed, confirm_secret, fio, role_id)
            -- VALUES (12, "password", "admin", 1,123,"fiiiioooorere",1);

create table locorg (
    `location_id` int not null, 
    `organizer_id` int not null
    );