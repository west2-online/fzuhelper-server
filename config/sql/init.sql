create table fzuhelper.`student`
(
    `id`                bigint              not null comment '学号',
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
    `created_at`        timestamp           default current_timestamp                   not null,
    `updated_at`        timestamp           default current_timestamp                   not null on update current_timestamp comment 'update profile time',
    `deleted_at`        timestamp           default null null,
    constraint `id`
        primary key (`id`)
)engine=InnoDB default charset=utf8mb4;

create table fzuhelper.`term`
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
            references fzuhelper.`student` (`id`)
            on delete cascade
)engine=InnoDB default charset=utf8mb4;

create table fzuhelper.`course`
(
    `id`                bigint              not null comment   '课程ID',
    `stu_id`            bigint              not null comment   '学生ID',
    `term_id`           bigint              not null comment   '学期ID',
    `type`              varchar(255)        not null comment   '修读类别',
    `name`              varchar(255)        not null comment   '课程名称',
    `paymentstatus`     varchar(255)        not null comment   '缴费状态',
    `syllabus`          varchar(255)        not null comment   '课程大纲',
    `lessonplan`        varchar(255)        not null comment   '课程计划',
    `credit`            decimal             not null comment   '学分',
    `electivetype`      varchar(255)        not null comment   '选课类型',
    `examtype`          varchar(255)        not null comment   '考试类别',
    `teacher`           varchar(255)        not null comment   '任课教师',
    `classroom`         varchar(255)        not null comment   '上课时间地点',
    `examtime`          varchar(255)        not null comment   '考试时间地点',
    `remark`            varchar(255)        not null comment   '备注',
    `adjust`            varchar(255)        not null comment   '调课信息',
    `created_at`        timestamp           default current_timestamp                   not null,
    `updated_at`        timestamp           default current_timestamp                   not null on update current_timestamp comment 'update profile time',
    `deleted_at`        timestamp           default null null,
    constraint `id`
        primary key (`id`),
    constraint `course_student`
        foreign key (`stu_id`)
            references fzuhelper.`student` (`id`)
            on delete cascade
)engine=InnoDB default charset=utf8mb4;

create table fzuhelper.`mark`
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
            references fzuhelper.`student` (`id`)
            on delete cascade
)engine=InnoDB default charset=utf8mb4;

create table fzuhelper.`picture` (
    `picture_id`        bigint              not null        comment '图片id',
    `url`               longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '图片链接',
    `href`              longtext            not null        comment '跳转链接',
    `text`              longtext                            comment '图片描述',
    `pic_type`          tinyint             default '1'     comment '图片类型',
    `show_times`        bigint              default '0'     comment '图片展示次数',
    `point_times`       bigint              default '0'     comment '图片点击次数',
    `duration`          bigint              default '3'     comment '图片展示持续时间',
    `start_at`          timestamp           default null    comment '展示开始时间',
    `end_at`            timestamp           default null    comment '展示结束时间',
    `start_time`        bigint              default '0'     comment '展示开始时间戳',
    `end_time`          bigint              default '24'    comment '展示结束时间戳',
    `s_type`            tinyint             default null    comment '图片类型',
    `frequency`         bigint              default null    comment '单天最大展示次数',
    `created_at`        timestamp           default current_timestamp   not null comment 'create time',
    `updated_at`        timestamp           default current_timestamp   not null on update current_timestamp comment 'update time',
    `deleted_at`        timestamp           default null    comment 'delete time',
  constraint `picture_id`
      primary key (`picture_id`),
  index `idx_picture_deleted_at` (`deleted_at`),
  index `retindex` (`start_at`,`end_at`,`start_time`,`end_time`,`s_type`)
) engine=InnoDB default charset=utf8mb4;