DROP DATABASE IF EXISTS video;
CREATE DATABASE video CHARACTER SET utf8 COLLATE utf8_general_ci;

USE video;

DROP TABLE IF EXISTS video;
CREATE TABLE video
(
  id        INT UNSIGNED UNIQUE NOT NULL AUTO_INCREMENT PRIMARY KEY,
  uid       CHAR(36) UNIQUE,
  title     VARCHAR(255)        NOT NULL,
  status    TINYINT UNSIGNED DEFAULT 1,
  duration  INT UNSIGNED                 DEFAULT 0,
  file_name VARCHAR(255)        NOT NULL
) Engine=InnoDB;