create table if not exists `user`
(
    id         int(10) auto_increment comment '主键id'
        primary key,
    sex        int default 0 not null comment '性别，0未知，1男2女',
    username   varchar(36)   not null comment '用户名',
    notice     varchar(36)   null comment '通知key',
    nickname   varchar(255)  null comment '昵称',
    email      varchar(255)  null comment '邮箱',
    phone      varchar(20)   null comment '手机号',
    avatar     varchar(255)  null comment '头像',
    login_time timestamp     null comment '最后登录时间',
    constraint user_pk2
        unique (username)
)
    comment '用户信息表';

