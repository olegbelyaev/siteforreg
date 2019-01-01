create database siteforeg;
use siteforeg;

create table locations (
    `id` int auto_increment primary key not null, 
    `name` varchar(255) not null default "", 
    `address` varchar(255) not null default ""
    );

create table orgaizesrs (
    `id` int auto_increment primary key not null, 
    `login` varchar(30) not null default "", 
    `password` varchar(30) not null default "", 
    `fio` varchar(250) not null default ""
    );

create table locorg (
    `location_id` int not null, 
    `organizer_id` int not null
    );