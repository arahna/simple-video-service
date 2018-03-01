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

INSERT INTO
  video
SET
  uid = 'd290f1ee-6c54-4b01-90e6-d701748f0851',
  title = 'Black Retrospective Woman',
  status = 3,
  duration = 127,
  file_name = 'index.mp4';
INSERT INTO
  video
SET
  uid = 'hjkhhjk3-23j4-j45k-erkj-kj3k4jl2k345',
  title = 'N Dance',
  status = 3,
  duration = 127,
  file_name = 'index.mp4';

INSERT INTO
  video
SET
  uid = 'sldjfl34-dfgj-523k-jk34-5jk3j45klj34',
  title = 'Cars',
  status = 3,
  duration = 127,
  file_name = '/content/sldjfl34-dfgj-523k-jk34-5jk3j45klj34/index.mp4';