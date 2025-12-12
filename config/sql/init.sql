create table `fzu-helper`.`student`
(
    `stu_id`            varchar(16)         not null                comment '学号',
    `name`              varchar(30)         not null                comment '姓名',
    `birthday`          varchar(12)         not null                comment '生日',
    `sex`               varchar(8)          not null                comment '性别',
    `college`           varchar(255)        not null                comment '学院',
    `grade`             bigint              not null                comment '年级',
    `major`             varchar(255)        not null                comment '专业',
    `created_at`        timestamp           not null default  current_timestamp,
    `updated_at`        timestamp           not null default  current_timestamp on update current_timestamp comment 'update profile time',
    `deleted_at`        timestamp           default  null null,
    constraint `id`
        primary key (`stu_id`),
    index `stu_birth`(`birthday`(10))
)engine=InnoDB default charset=utf8mb4;

create table `fzu-helper`.`term`
(
    `id`                bigint              not null comment '学期ID',
    `stu_id`            bigint              not null comment '学生ID',
    `term_time`         varchar(255)        not null comment '学期时间',
    `created_at`        timestamp           default current_timestamp                   not null,
    `updated_at`        timestamp           default current_timestamp                   not null on update current_timestamp comment 'update profile time',
    `deleted_at`        timestamp           default null null,
    constraint `id`
        primary key (`id`)
)engine=InnoDB default charset=utf8mb4;

CREATE TABLE `fzu-helper`.`scores` (
                                       `stu_id` varchar(16) NOT NULL COMMENT '学生ID',
                                       `scores_info` json NOT NULL COMMENT '学生成绩信息',
                                       `scores_info_sha256` varchar(64) NOT NULL COMMENT '学生成绩信息SHA256',
                                       `created_at` timestamp NOT NULL DEFAULT current_timestamp,
                                       `updated_at` timestamp NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp,
                                       `deleted_at` timestamp NULL DEFAULT NULL,
                                       PRIMARY KEY (`stu_id`)
) ENGINE = InnoDB CHARSET = utf8mb4;

CREATE TABLE `fzu-helper`.`course_offerings` (
                                                 `id` BIGINT NOT NULL AUTO_INCREMENT,
                                                 `name` VARCHAR(64) NOT NULL COMMENT '课程名',
                                                 `term` VARCHAR(16) NOT NULL COMMENT '学期',
                                                 `teacher` VARCHAR(255) NOT NULL COMMENT '教师全名',
                                                 `elective_type` VARCHAR(64) NOT NULL COMMENT '选修类型',
                                                 `course_hash` CHAR(64)  NOT NULL COMMENT '通过name、term、teacher、elective_type生成的唯一hash',
                                                 `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                                 `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                                 `deleted_at` TIMESTAMP NULL DEFAULT NULL,
                                                 PRIMARY KEY (`id`),
                                                 UNIQUE INDEX `uniq_course_hash` (`course_hash`)
) ENGINE=InnoDB CHARSET=utf8mb4;

create table `fzu-helper`.`launch_screen`(
    `id`          bigint                NOT NULL           AUTO_INCREMENT           COMMENT 'ID',
    `url`         tinytext              NULL                                        COMMENT '图片url',
    `href`        tinytext              NULL                                        COMMENT '示例:"Toapp:abab"',
    `text`        varchar(255)          NULL                                        COMMENT '图片描述',
    `pic_type`    bigint                NOT NULL           DEFAULT 1                COMMENT '1为空，2为页面跳转，3为app跳转',
    `show_times`  bigint                NOT NULL           DEFAULT 0                COMMENT '展示次数(GetMobileImage)',
    `point_times` bigint                NOT NULL           DEFAULT 0                COMMENT '点击次数(AddPointTime)',
    `duration`    bigint                NOT NULL           DEFAULT 3                COMMENT '展示时间，直接从客户端传入的值',
    `start_at`    timestamp             NULL                                        COMMENT '开始时间',
    `end_at`      timestamp             NULL                                        COMMENT '结束时间',
    `start_time`  bigint                NOT NULL           DEFAULT 0                COMMENT '开始时段 0-24',
    `end_time`    bigint                NOT NULL           DEFAULT 24               COMMENT '结束时段 0-24',
    `s_type`      bigint                NULL                                        COMMENT '类型',
    `frequency`   bigint                NULL                                        COMMENT '一天展示频率',
    `regex`       mediumtext            NOT NULL                                    COMMENT '存储所有投放学生学号的json',
    `created_at` timestamp              NOT NULL           DEFAULT current_timestamp ,
    `updated_at` timestamp              NOT NULL           DEFAULT current_timestamp ON UPDATE current_timestamp,
    `deleted_at` timestamp              NULL               DEFAULT NULL,
    constraint `id`
        primary key (`id`)
)ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `fzu-helper`.`course`(
    `id`                  bigint      NOT NULL COMMENT 'ID',
    `stu_id`              varchar(16) NOT NULL COMMENT '学生ID',
    `term`                varchar(16)  NOT NULL COMMENT '学期',
    `term_courses`        json        NOT NULL COMMENT '学期课程信息',
    `term_courses_sha256` varchar(64) NOT NULL COMMENT '学期课程信息SHA256',
    `created_at`          timestamp   NOT NULL DEFAULT current_timestamp,
    `updated_at`          timestamp   NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp,
    `deleted_at`          timestamp   NULL     DEFAULT NULL,
    key `term` (`term`),
    constraint `id`
        primary key (`id`)
)engine=InnoDB default charset=utf8mb4;

CREATE TABLE `fzu-helper`.`notice`(
    `id`          bigint      NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `title`       varchar(255) NOT NULL COMMENT '标题',
    `url`         varchar(255)         NOT NULL COMMENT '链接',
    `published_at` varchar(10)    NOT NULL COMMENT '发布时间',
    `created_at`  timestamp    NOT NULL DEFAULT current_timestamp,
    `updated_at`  timestamp    NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp,
    `deleted_at`  timestamp    NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `unique_url` UNIQUE (`url`)
)engine=InnoDB default charset=utf8mb4;
/* 建立发布时间的索引 */
CREATE INDEX idx_published_at ON `fzu-helper`.`notice`(`published_at`);

CREATE TABLE `fzu-helper`.`visit`(
    `id`          bigint       NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `date`         varchar(12)  NOT NULL                COMMENT '日期',
    `visits`       bigint       NOT NULL                COMMENT '访问统计',
    `created_at`  timestamp    NOT NULL DEFAULT current_timestamp,
    `updated_at`  timestamp    NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp,
    `deleted_at`  timestamp        NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
#     index visit_date(`date`), ## UNIQUE会隐式的建立一个索引
    CONSTRAINT `unique_date` UNIQUE (`date`)
)engine=InnoDB default charset=utf8mb4;

CREATE TABLE `fzu-helper`.`toolbox_config` (
    `id`          bigint       NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `tool_id`     bigint       NOT NULL COMMENT '工具ID',
    `visible`     tinyint      NOT NULL DEFAULT 1 COMMENT '是否可见',
    `name`        varchar(255) COMMENT '功能名称',
    `icon`        varchar(255) COMMENT '图标网址',
    `type`        varchar(255) COMMENT '工具类型',
    `message`     varchar(255) COMMENT '消息',
    `extra`       varchar(255) COMMENT '额外信息',
    `student_id`  varchar(255) COMMENT '学号白名单',
    `platform`    varchar(255) COMMENT 'android或ios白名单',
    `version`     bigint       COMMENT '版本号',
    `created_at`  timestamp    NOT NULL DEFAULT current_timestamp,
    `updated_at`  timestamp    NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp,
    `deleted_at`  timestamp    NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_toolbox_config` (`tool_id`, `student_id`, `platform`, `version`)
) engine=InnoDB default charset=utf8mb4;

CREATE TABLE IF NOT EXISTS `admin_secrets` (
    `id`          bigint       NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `module_name` varchar(255) NOT NULL COMMENT '模块名称，如：toolbox, notice, user等',
    `secret_key`  varchar(255) NOT NULL COMMENT '密钥值',
    `created_at`  timestamp    NOT NULL DEFAULT current_timestamp COMMENT '创建时间',
    `updated_at`  timestamp    NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp COMMENT '更新时间',
    `deleted_at`  timestamp    NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_module_secret` (`module_name`, `secret_key`)
) engine=InnoDB default charset=utf8mb4;


CREATE TABLE `fzu-helper`.`follow_relation`
(
    `id`           bigint        NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `follower_id`  varchar(16)   NOT NULL COMMENT '关注者学号',
    `followed_id`  varchar(16)   NOT NULL COMMENT '被关注者学号',
    `status`       tinyint       NOT NULL DEFAULT 0 COMMENT '状态: 0-关注中, 1-已取消关注',
    `created_at`   timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`   timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at`   timestamp     NULL DEFAULT NULL,
    CONSTRAINT `pk_id` PRIMARY KEY (`id`),
    UNIQUE KEY `uk_follower_followed` (`follower_id`, `followed_id`),
    INDEX `idx_follower_id` (`follower_id`),
    INDEX `idx_followed_id` (`followed_id`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='关注关系表';
