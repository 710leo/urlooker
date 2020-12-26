CREATE DATABASE `urlooker`;
USE urlooker;

#DROP TABLE IF EXISTS `strategy`;
CREATE TABLE `strategy` ( 
  `id`          int(10)         unsigned NOT NULL AUTO_INCREMENT,
  `url`         varchar(1024)   NOT NULL,
  `idc`         varchar(255)    NOT NULL DEFAULT '',
  `enable`      int(1)          NOT NULL DEFAULT 1,
  `keywords`    varchar(255)    NOT NULL DEFAULT '',
  `endpoint`    varchar(255)    NOT NULL DEFAULT '',
  `timeout`     varchar(255)    NOT NULL DEFAULT '',
  `creator`     varchar(255)    NOT NULL DEFAULT '',
  `data`        text,
  `ip`          varchar(255)    NOT NULL DEFAULT '', 
  `expect_code` varchar(255)    NOT NULL DEFAULT '', 
  `tag`         varchar(255)    NOT NULL DEFAULT '',
  `header`      text,
  `method`      varchar(255)    DEFAULT "get",
  `post_data`   text,
  `note`        text,
  `ding_webhook` varchar(255)   NOT NULL DEFAULT '',
  `max_step`    int(4)          NOT NULL DEFAULT 3,
  `times`       int(4)          NOT NULL DEFAULT 3,
  `teams`       varchar(32)     NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci; 
#ALTER TABLE strategy ADD ip varchar(255);
#ALTER TABLE strategy ADD enable int(1) DEFAULT 1;
#ALTER TABLE strategy ADD endpoint varchar(255) DEFAULT "";
#ALTER TABLE strategy ADD method varchar(255) DEFAULT "";
#ALTER TABLE strategy ADD header text;
#ALTER TABLE strategy ADD post_data text;
#ALTER TABLE strategy ADD idc varchar(255);

#DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id`        BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT,
  `name`      VARCHAR(64)      NOT NULL,
  `cnname`    VARCHAR(64)      NOT NULL DEFAULT '',
  `password`  VARCHAR(32)      NOT NULL,
  `email`     VARCHAR(255)     NOT NULL DEFAULT '',
  `phone`     VARCHAR(16)      NOT NULL DEFAULT '',
  `wechat`    VARCHAR(255)     NOT NULL DEFAULT '',
  `role`      TINYINT          NOT NULL DEFAULT 0,
  `created`   TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

#DROP TABLE if EXISTS team;
CREATE TABLE `team` (
  `id`      BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT,
  `name`    VARCHAR(64)      NOT NULL,
  `resume`  VARCHAR(255)     NOT NULL DEFAULT '',
  `creator` BIGINT UNSIGNED  NOT NULL,
  `created` TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_team_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

#DROP TABLE if EXISTS `rel_team_user`;
CREATE TABLE `rel_team_user` (
  `id`    BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `tid`   BIGINT UNSIGNED NOT NULL,
  `uid`   BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_rel_tid` (`tid`),
  KEY `idx_rel_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

#DROP TABLE IF EXISTS `rel_sid_ip`;
CREATE TABLE `rel_sid_ip` (
  `id`        INT UNSIGNED  NOT NULL AUTO_INCREMENT,
  `sid`       INT UNSIGNED  NOT NULL,
  `ip`        VARCHAR(32)   NOT NULL DEFAULT '',
  `ts`        INT(10),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_sid_ip` (`sid`,`ip`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

#DROP TABLE IF EXISTS `event`;
CREATE TABLE `event` (
  `id`          INT UNSIGNED    NOT NULL AUTO_INCREMENT,
  `event_id`    VARCHAR(64)     NOT NULL,
  `status`      VARCHAR(32)     NOT NULL,
  `url`         VARCHAR(256)    NOT NULL DEFAULT '',
  `ip`          VARCHAR(32)     NOT NULL DEFAULT '',
  `strategy_id` INT,
  `event_time`  INT(11),
  `resp_time`   INT(6),
  `resp_code`   VARCHAR(3),
  `result`      INT(1)          NOT NULL DEFAULT 0 COMMENT '0:no error, 1:timeout, 2:expect code err, 3,keyword unmatch 4:dns err', 
  `current_step`INT(1),
  `max_step`    INT(1),
  PRIMARY KEY (`id`),
  INDEX `idx_strategy_id` (`strategy_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

#DROP TABLE IF EXISTS `item_status00`;
CREATE TABLE `item_status00` (
  `id`        INT UNSIGNED  NOT NULL AUTO_INCREMENT,
  `sid`       INT UNSIGNED  NOT NULL,
  `ip`        VARCHAR(32)   NOT NULL DEFAULT '',
  `resp_time` INT(6),
  `resp_code` VARCHAR(3),
  `push_time` INT(10),
  `result`    INT(1)        NOT NULL DEFAULT 0 COMMENT '0:no error, 1:timeout, 2:expect code err, 3,keyword unmatch 4:dns err', 
  PRIMARY KEY (`id`),
  INDEX `idx_ip` (`ip`),
  INDEX `idx_sid` (`sid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

