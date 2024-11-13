create table `fzu-helper`.`student`
(
    `id`                bigint              not null comment '学生ID',
    `number`            varchar(255)        not null comment '学号',
    `password`          varchar(255)        not null comment '密码',
    `sex`               varchar(255)        not null comment '性别',
    `birthday`          varchar(255)        not null comment '出生日期',
    `phone`             varchar(255)        not null comment '手机号',
    `email`             varchar(255)                 comment '邮箱',
    `college`           varchar(255)        not null comment '学院',
    `grade`             bigint              not null comment '年级',
    `status_change`     varchar(255)                 comment '学籍异动与奖励',
    `major`             varchar(255)        not null comment '专业',
    `counselor`         varchar(255)        not null comment '辅导员',
    `examinee_category` varchar(255)        not null comment '考生类别',
    `nationality`       varchar(255)        not null comment '民族',
    `country`           varchar(255)        not null comment '国别',
    `political_status`  varchar(255)        not null comment '政治面貌',
    `source`            varchar(255)        not null comment '生源地',
    `created_at`        timestamp           default  current_timestamp                   not null,
    `updated_at`        timestamp           default  current_timestamp                   not null on update current_timestamp comment 'update profile time',
    `deleted_at`        timestamp           default  null null,
    constraint `id`
        primary key (`id`)
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
        primary key (`id`),
    constraint `term_student`
        foreign key (`stu_id`)
            references `fzu-helper`.`student` (`id`)
            on delete cascade
)engine=InnoDB default charset=utf8mb4;

create table `fzu-helper`.`mark`
(
    `id`                bigint              not null comment   '成绩ID',
    `stu_id`            bigint              not null comment   '学生ID',
    `term_id`           bigint              not null comment   '学期ID',
    `course_id`         bigint              not null comment   '课程ID',
    `type`              varchar(255)        not null comment   '修读类别',
    `semester`          varchar(255)        not null comment   '开课学期',
    `name`              varchar(255)        not null comment   '课程名称',
    `credit`            decimal             not null comment   '计划学分',
    `score`             varchar(255)        not null comment   '得分',
    `gpa`               varchar(255)        not null comment   '绩点',
    `earned_credits`    decimal             not null comment   '得到学分',
    `electivetype`      varchar(255)        not null comment   '选课类型',
    `examtype`          varchar(255)        not null comment   '考试类别',
    `teacher`           varchar(255)        not null comment   '任课教师',
    `classroom`         varchar(255)        not null comment   '上课时间地点',
    `examtime`          varchar(255)        not null comment   '考试时间地点',
    `created_at`        timestamp           default current_timestamp                   not null,
    `updated_at`        timestamp           default current_timestamp                   not null on update current_timestamp comment 'update profile time',
    `deleted_at`        timestamp           default null null,
    constraint `id`
        primary key (`id`),
    constraint `mark_student`
        foreign key (`stu_id`)
            references `fzu-helper`.`student` (`id`)
            on delete cascade
)engine=InnoDB default charset=utf8mb4;

create table `fzu-helper`.`launch_screen`(
    `id`          bigint                NOT NULL           AUTO_INCREMENT           COMMENT 'ID',
    `url`         varchar(512)          NULL                                        COMMENT '图片url',
    `href`        varchar(255)          NULL                                        COMMENT '示例:"Toapp:abab"',
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
    `term`                varchar(8)  NOT NULL COMMENT '学期',
    `term_courses`        json        NOT NULL COMMENT '学期课程信息',
    `term_courses_sha256` varchar(64) NOT NULL COMMENT '学期课程信息SHA256',
    `created_at`          timestamp   NOT NULL DEFAULT current_timestamp,
    `updated_at`          timestamp   NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp,
    `deleted_at`          timestamp   NULL     DEFAULT NULL,
    key `term` (`term`),
    constraint `id`
        primary key (`id`)
)
