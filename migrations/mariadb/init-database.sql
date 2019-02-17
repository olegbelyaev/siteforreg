create database IF NOT EXISTS siteforeg;
use siteforeg;

create table IF NOT EXISTS locations (
    `id` int(11) auto_increment primary key not null, 
    `name` varchar(255) not null default "", 
    `address` varchar(255) not null default ""
    );


CREATE TABLE IF NOT EXISTS `users` (
    `id` int(11) NOT NULL AUTO_INCREMENT primary key , 
    `password` varchar(30) NOT NULL DEFAULT "", 
    `email` varchar(100) NOT NULL DEFAULT "", 
    `fio` varchar(250)  NOT NULL DEFAULT "", 
    `roles` Set('listener','organizer','admin')
    );


create table IF NOT EXISTS locorg (
    `location_id` int not null, 
    `organizer_id` int not null
    );


CREATE TABLE IF NOT EXISTS `lectures` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `location_id` int(11) NOT NULL DEFAULT 0,
  `when` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `group_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `max_seets` int(11) NOT NULL DEFAULT 30,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `description` varchar(1024) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`));

CREATE TABLE IF NOT EXISTS `tickets` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`user_id` int(11) NOT NULL,
`lecture_id` int(11) NOT NULL,
PRIMARY KEY (`id`)
);