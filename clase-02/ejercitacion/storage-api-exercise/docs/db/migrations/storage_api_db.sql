-- DDL: Data Definition Language
DROP DATABASE IF EXISTS `storage_api_db`;

CREATE DATABASE `storage_api_db`;

USE `storage_api_db`;

CREATE TABLE `products` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `quantity` int NOT NULL,
  `code_value` varchar(255) NOT NULL,
  `is_published` boolean NOT NULL,
  `expiration` date NOT NULL,
  `price` decimal(10, 2) NOT NULL,
  PRIMARY KEY (`id`)
);