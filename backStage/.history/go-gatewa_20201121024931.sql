/*
    Host:127.0.0.1
    Generation:2020-11-08
    server edition ： 5.7.32
*/

-- --------------------------------------------------------------------------------|
-- --------------------------------------------------------------------------------|
-- --------------------------------set database------------------------------------|
/*
    创建数据库(go_gateway)
*/
create database if not exists go_gateway
default character set utf8
default collate utf8_general_ci;

/* !drop database if exists go_gateway*/
/* !SELECT concat('DROP TABLE IF EXISTS ', table_name, ';') 
    FROM information_schema.tables 
    WHERE table_schema = 'go_gateway' 检索出所有可以删除的表的语句 */
-- -------------------------------------------------------------------

/* 影响AUTO_INCREMETNT列的处理，
一般情况下你可以向该列插入NULL或0生成下一个序列号。
但是NO_AUTO_VALUE_ON_ZERO禁用0，
因此只有NULL可以生成下一个序列号 */
set SQL_MODE = `NO_AUTO_VALUE_ON_ZERO`; 

/* 当前session禁用自动提交事务，
自此句执行以后，每个SQL语句或者语句块所在的事务都需要显示commit才能提交事务 */
SET AUTOCOMMIT = 0; 

/* 启动一个新事务 */
START TRANSACTION; 
/* commit */

/* 修改当前会话时区 */
SET time_zone=`+00:00`; 

-- ---------------------------------block end--------------------------------------|
-- --------------------------------------------------------------------------------|
-- --------------------------------------------------------------------------------|

-- --------------------------------------------------------------------------------|
-- --------------------------------------------------------------------------------|
-- ---------------------------------set tables-------------------------------------|
-- 
-- Database : `go_gateway`
--

--

--
-- 表的结构：`gateway_admin`
--

CREATE TABLE if not exists go_gateway.gateway_admin  (
    `id`        bigint(20)      not null                                    comment '自增id',
    `user_name` varchar(255)    not null default ''                         comment '用户名',
    `salt`      varchar(50)     not null default ''                         comment '盐',
    `password`  varchar(255)    not null default ''                         comment '密码',
    `create_at` datetime        not null default '1971-01-01 00:00:00'      comment '新增时间' ,
    `update_at` datetime        not null default '1971-01-01 00:00:00'      comment '更新时间',
    `is_delete` tinyint(4)      not null default '0'                        comment '是否删除'
) engine=InnoDB default charset=utf8 comment='管理员表';

/* ! drop table if exists go_gateway.gateway_admin 测试用语句 */

--
-- 向gateway_admin中存入测试用数据
--

insert into 
go_gateway.gateway_admin (`id`, `user_name`, `salt`, `password`, `create_at`, `update_at`, `is_delete`)
values 
(1, 'admin', 'admin', '2823d896e9822c0833d41d4904f0c00756d718570fce49b9a379a62c804689d3', (SELECT DATE_FORMAT(NOW(),'%Y-%m-%d %H:%i:%s')), (SELECT DATE_FORMAT(NOW(),'%Y-%m-%d %H:%i:%s')), 0);
/* 
delete from go_gateway.gateway_admin 清除表数据
select * from go_gateway.gateway_admin 查询数据是否成功添加
*/

-- -------------------------------------------------------------------
--
-- 表的结构 `gateway_app`
--

CREATE TABLE if not exists go_gateway.gateway_app (
    /* bigint 一个略大整数，-2^63 (-9223372036854775808) 到 2^63-1 (9223372036854775807) 的整型数据。【远大于21亿左右】
       bigint(20)	存储数字：666	存储空间：8	实际显示宽度：20	实际显示：00000000000000000666*/
    /* unsigned无符号化 */
    /* default 设置默认值 */
    /* tinyint 一个极小整数
       tinyint(3)	存储数字：33	存储空间：1	实际显示宽度：3	实际显示：033*/
    `id`            bigint(20)      unsigned not null               comment '自增id',
    `app_id`        varchar(255)             not null default ''    comment '租户id',
    `name`          varchar(255)             not null default ''    comment '租户名称',
    `secret`        varchar(255)             not null default ''    comment '密钥',
    `white_ips`     varchar(1000)            not null default ''    comment 'ip白名单，支持前缀匹配',
    `qpd`           bigint(20)               NOT NULL default '0'   comment '日请求量限制',
    `qps`           bigint(20)               not null default '0'   comment '每秒请求量限制',
    `create_at`     datetime                 not null               comment '添加时间',
    `update_at`     datetime                 not null               comment '更新时间',
    `is_delete`     tinyint(4)               not null default '0'   comment '是否删除 1=删除'
) ENGINE=InnoDB default charset=utf8 comment='网关租户表';

/* !drop table if exists go_gateway.gateway_app 测试用语句 */

--
-- 向go_gateway.gateway_app表中添加测试数据
--

insert into go_gateway.gateway_app (id, app_id, name, secret, white_ips, qpd, qps, create_at, update_at, is_delete)
values (31, 'app_id_a', '租户A', '449441eb5e72dca9c42a12f3924ea3a2', 'white_ips', 100000, 100, (SELECT DATE_FORMAT(NOW(),'%Y-%m-%d %H:%i:%s')),(SELECT DATE_FORMAT(NOW(),'%Y-%m-%d %H:%i:%s')), 0),
(32, 'app_id_b', '租户B', '8d7b11ec9be0e59a36b52f32366c09cb', '', 20, 0, (SELECT DATE_FORMAT(NOW(),'%Y-%m-%d %H:%i:%s')), (SELECT DATE_FORMAT(NOW(),'%Y-%m-%d %H:%i:%s')), 0),
(33, 'app_id', '租户名称', '', '', 0, 0, (SELECT DATE_FORMAT(NOW(),'%Y-%m-%d %H:%i:%s')), (SELECT DATE_FORMAT(NOW(),'%Y-%m-%d %H:%i:%s')), 1),
(34, 'app_id45', '名称', '07d980f8a49347523ee1d5c1c41aec02', '', 0, 0, (SELECT DATE_FORMAT(NOW(),'%Y-%m-%d %H:%i:%s')), (SELECT DATE_FORMAT(NOW(),'%Y-%m-%d %H:%i:%s')), 1);

/*
    delete from go_gateway.gateway_app 清除测试数据
    select * from go_gateway.gateway_app 查询测试数据是否添加成功
*/

-- -------------------------------------------------------------------
--
-- 表的结构 `gateway_service_access_control`
--

create table if not exists go_gateway.gateway_service_access_control (
    `id`                    bigint(20)          not null                comment '自增主键',
    `service_id`            bigint(20)          not null default '0'    comment '服务id',
    `open_auth`             tinyint(4)          not null default '0'    comment '是否开启权限 1=开启',
    `black_list`            varchar(1000)       not null default ''     comment '黑名单ip',
    `white_list`            varchar(1000)       not null default ''     comment '白名单ip',
    `white_host_home`       varchar(1000)       not null default ''     comment '白名单主机',
    `clientip_flow_limit`   int(11)             not null default '0'    comment '客户端ip限流',
    `service_flow_limit`    int(20)             not null default '0'    comment '服务器限流'
)engine=InnoDB default charset=utf8 comment='网关权限控制表';

/* ! drop table if exists go_gateway.gateway_service_access_control */

--
-- 向go_gateway.gateway_service_access_controls中添加测试数据
--

insert into go_gateway.gateway_service_access_control
(id,service_id,open_auth,black_list,white_list,white_host_home,clientip_flow_limit,service_flow_limit)
values
(162, 35, 1, '', '', '', 0, 0),
(165, 34, 0, '', '', '', 0, 0),
(167, 36, 0, '', '', '', 0, 0),
(168, 38, 1, '111.11', '22.33', '11.11', 12, 12),
(169, 41, 1, '111.11', '22.33', '11.11', 12, 12),
(170, 42, 1, '111.11', '22.33', '11.11', 12, 12),
(171, 43, 0, '111.11', '22.33', '11.11', 12, 12),
(172, 44, 0, '', '', '', 0, 0),
(173, 45, 0, '', '', '', 0, 0),
(174, 46, 0, '', '', '', 0, 0),
(175, 47, 0, '', '', '', 0, 0),
(176, 48, 0, '', '', '', 0, 0),
(177, 49, 0, '', '', '', 0, 0),
(178, 50, 0, '', '', '', 0, 0),
(179, 51, 0, '', '', '', 0, 0),
(180, 52, 0, '', '', '', 0, 0),
(181, 53, 0, '', '', '', 0, 0),
(182, 54, 1, '127.0.0.3', '127.0.0.2', '', 11, 12),
(183, 55, 1, '127.0.0.2', '127.0.0.1', '', 45, 34),
(184, 56, 0, '192.168.1.0', '', '', 0, 0),
(185, 57, 0, '', '127.0.0.1,127.0.0.2', '', 0, 0),
(186, 58, 1, '', '', '', 0, 0),
(187, 59, 1, '127.0.0.1', '', '', 2, 3),
(188, 60, 1, '', '', '', 0, 0),
(189, 61, 0, '', '', '', 0, 0);

/*
    select * from go_gateway.gateway_service_access_control 查询测试数据是否添加成功
    delete from go_gateway.gateway_service_access_control 删除测试数据
*/

-- -------------------------------------------------------------------
--
-- 表的结构 `go_gateway.gateway_service_grpc_rule`
--

create table if not exists go_gateway.gateway_service_grpc_rule (
    `id`                bigint(20)      not null                comment '自增主键',
    `service_id`        bigint(20)      not null default '0'    comment '服务id',
    `port`              int(5)          not null default '0'    comment '端口',
    `header_transfor`   varchar(5000)   not null default ''     comment 'header转换支持增加(add)、删除(del)、修改(edit) 格式：add headname headvalue 多个逗号间隔'
)engine=InnoDB default charset=utf8 comment='网关路由匹配表';

--
-- 向go_gateway.gateway_service_grpc_rule添加测试数据
--

insert into go_gateway.gateway_service_grpc_rule 
(id, service_id, port, header_transfor)
values 
(171, 53, 8009, ''),
(172, 54, 8002, 'add metadata1 datavalue,edit metadata2 datavalue2'),
(173, 58, 8012, 'add meta_name meta_value');

/*
    select * from go_gateway.gateway_service_grpc_rule 查询测试数据是否添加成功
    delete from go_gateway.gateway_service_grpc_rule 删除测试数据
*/

-- -------------------------------------------------------------------
--
-- 表的结构 `go_gateway.gateway_service_http_rule`
--

create table if not exists go_gateway.gateway_service_http_rule (
    `id`             bigint(20)     not null                    comment '自增主键',
    `service_id`     bigint(20)     not null                    comment '服务id',
    `rule_type`      tinyint(4)     not null default '0'        comment '匹配类型 0=url前缀url_prefix 1=域名domain',
    `rule`           varchar(255)   not null default ''         comment 'type=domain表示域名，type=url_prefix时标是url前缀',
    `need_https`     tinyint(4)     not null default '0'        comment '支持https 1=支持',
    `need_strip_url` tinyint(4)     not null default '0'        comment '启用strip_url 1=启用',
    `need_websocket` tinyint(4)     not null default '0'        comment '是否支持websocket 1=支持',
    `url_rewrite`    varchar(5000)  not null default ''         comment 'url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔',
    `header_transfor`varchar(5000)  not null default ''         comment 'header转换支持增加(add)、删除(del)、修改(edit) 格式：add headname headvalue 多个逗号间隔'
)engine=InnoDB default charset=utf8 comment='网关路由匹配表';

/* ! drop table if exists go_gateway.gateway_service_http_rule */

--
-- 向go_gateway.gateway_service_http_rule添加测试数据
--

insert into go_gateway.gateway_service_http_rule 
(id,service_id,rule_type,rule,need_https,need_strip_url,need_websocket,url_rewrite,header_transfor)
values
(165, 35, 1, '', 0, 0, 0, '', ''),
(168, 34, 0, '', 0, 0, 0, '', ''),
(170, 36, 0, '', 0, 0, 0, '', ''),
(171, 38, 0, '/abc', 1, 0, 1, '^/abc $1', 'add head1 value1'),
(172, 43, 0, '/usr', 1, 1, 0, '^/afsaasf $1,^/afsaasf $1', ''),
(173, 44, 1, 'www.test.com', 1, 1, 1, '', ''),
(174, 47, 1, 'www.test.com', 1, 1, 1, '', ''),
(175, 48, 1, 'www.test.com', 1, 1, 1, '', ''),
(176, 49, 1, 'www.test.com', 1, 1, 1, '', ''),
(177, 56, 0, '/test_http_service', 1, 1, 1, '^/test_http_service/abb/(.*) /test_http_service/bba/$1', 'add header_name header_value'),
(178, 59, 1, 'test.com', 0, 1, 1, '', 'add headername headervalue'),
(179, 60, 0, '/test_strip_uri', 0, 1, 0, '^/aaa/(.*) /bbb/$1', ''),
(180, 61, 0, '/test_https_server', 1, 1, 0, '', '');

/*
    select * from go_gateway.gateway_service_http_rule 查询测试数据是否添加成功
    delete from go_gateway.gateway_service_http_rule   删除测试数据
*/

-- -------------------------------------------------------------------
--
-- 表的结构 `go_gateway.gateway_service_info`
--

create table if not exists go_gateway.gateway_service_info (
    `id`                bigint(20)      unsigned     not null comment '自增主键',
    `load_type`         tinyint(4)                   not null default '0'                   comment '负载类型 0=http 1=tcp 2=grpc',
    `service_name`      varchar(255)                 not null default ''                    comment '服务名称 6-128 数字字母下划线',
    `service_desc`      varchar(255)                 not null default ''                    comment '服务描述',
    `create_at`         datetime                     not null default '1971-01-01 00:00:00' comment '添加时间',
    `update_at`         datetime                     not null default '1971-01-01 00:00:00' comment '更新时间',
    `is_delete`         tinyint(4)                            default '0'                   comment '是否删除 1=删除'
)engine=InnoDB default charset=utf8 comment='网关基本信息表';

/* !drop table if exists go_gateway.gateway_service_info */

--
-- 向go_gateway.gateway_service_info中添加测试数据
--

insert into go_gateway.gateway_service_info 
(id,load_type,service_name,service_desc,create_at,update_at,is_delete)
values
(34, 0, 'websocket_test', 'websocket_test', '2020-04-13 01:31:47', '1971-01-01 00:00:00', 1),
(35, 1, 'test_grpc', 'test_grpc', '2020-04-13 01:34:32', '1971-01-01 00:00:00', 1),
(36, 2, 'test_httpe', 'test_httpe', '2020-04-11 21:12:48', '1971-01-01 00:00:00', 1),
(38, 0, 'service_name', '11111', '2020-04-15 07:49:45', '2020-04-11 23:59:39', 1),
(41, 0, 'service_name_tcp', '11111', '2020-04-13 01:38:01', '2020-04-12 01:06:09', 1),
(42, 0, 'service_name_tcp2', '11111', '2020-04-13 01:38:06', '2020-04-12 01:13:24', 1),
(43, 1, 'service_name_tcp4', 'service_name_tcp4', '2020-04-15 07:49:44', '2020-04-12 01:13:50', 1),
(44, 0, 'websocket_service', 'websocket_service', '2020-04-15 07:49:43', '2020-04-13 01:20:08', 1),
(45, 1, 'tcp_service', 'tcp_desc', '2020-04-15 07:49:41', '2020-04-13 01:46:27', 1),
(46, 1, 'grpc_service', 'grpc_desc', '2020-04-13 01:54:12', '2020-04-13 01:53:14', 1),
(47, 0, 'testsefsafs', 'werrqrr', '2020-04-13 01:59:14', '2020-04-13 01:57:49', 1),
(48, 0, 'testsefsafs1', 'werrqrr', '2020-04-13 01:59:11', '2020-04-13 01:58:14', 1),
(49, 0, 'testsefsafs1222', 'werrqrr', '2020-04-13 01:59:08', '2020-04-13 01:58:23', 1),
(50, 2, 'grpc_service_name', 'grpc_service_desc', '2020-04-15 07:49:40', '2020-04-13 02:01:00', 1),
(51, 2, 'gresafsf', 'wesfsf', '2020-04-15 07:49:39', '2020-04-13 02:01:57', 1),
(52, 2, 'gresafsf11', 'wesfsf', '2020-04-13 02:03:41', '2020-04-13 02:02:55', 1),
(53, 2, 'tewrqrw111', '123313', '2020-04-13 02:03:38', '2020-04-13 02:03:20', 1),
(54, 2, 'test_grpc_service1', 'test_grpc_service1', '2020-04-15 07:49:37', '2020-04-15 07:38:43', 1),
(55, 1, 'test_tcp_service1', 'redis服务代理', '2020-04-15 07:49:35', '2020-04-15 07:46:35', 1),
(56, 0, 'test_http_service', '测试HTTP代理', '2020-04-16 00:54:45', '2020-04-15 07:55:07', 0),
(57, 1, 'test_tcp_service', '测试TCP代理', '2020-04-19 14:03:09', '2020-04-15 07:58:39', 0),
(58, 2, 'test_grpc_service', '测试GRPC服务', '2020-04-21 07:20:16', '2020-04-15 07:59:46', 0),
(59, 0, 'test.com:8080', '测试域名接入', '2020-04-18 22:54:14', '2020-04-18 20:29:13', 0),
(60, 0, 'test_strip_uri', '测试路径接入', '2020-04-21 06:55:26', '2020-04-18 22:56:37', 0),
(61, 0, 'test_https_server', '测试https服务', '2020-04-19 12:22:33', '2020-04-19 12:17:04', 0);

/*
    select * from go_gateway.gateway_service_info
    delete from go_gateway.gateway_service_info
*/

-- -------------------------------------------------------------------
--
-- 表的结构 `go_gateway.gateway_service_load_balance`
--

create table if not exists go_gateway.gateway_service_load_balance (
    `id`                        bigint(20)      not null                comment '自增主键',
    `service_id`                bigint(20)      not null default '0'    comment '服务id',
    `check_method`              tinyint(20)     not null default '0'    comment '检查方法 0=tcpchk，检测端口是否握手成功',
    `check_timeout`             int(10)         not null default '0'    comment 'check超时时间，单位s',
    `check_interval`            int(11)         not null default '0'    comment '检查间隔，单位s',
    `round_type`                tinyint(4)      not null default '2'    comment '轮询方式 0=random 1=round_robin 2=wieght_round_robin 3=ip_hash',
    `ip_list`                   varchar(2000)   not null default ''     comment 'ip列表',
    `weight_list`               varchar(2000)   not null default ''     comment '权重列表',
    `forbid_list`               varchar(2000)   not null default ''     comment '禁用ip列表',
    `upstream_connect_timeout`  int(11)         not null default '0'    comment '建立连接超时，单位s',
    `upstream_header_timeout`   int(11)         not null default '0'    comment '获取header超时，单位s',
    `upstream_idle_timeout`     int(10)         not null default '0'    comment '链接最大空闲时间，单位s',
    `upstream_max_idle`         int(11)         not null default '0'    comment '最大空闲链接数'
)engine=InnoDB default charset=utf8 comment='网关负载表';

/* !drop table if exists go_gateway.gateway_service_load_balance */

--
-- 向go_gateway.gateway_service_load_balance添加测试数据
--

insert into go_gateway.gateway_service_load_balance 
(id,service_id,check_method,check_timeout,check_interval,round_type,ip_list,weight_list,forbid_list,upstream_connect_timeout,upstream_header_timeout,upstream_idle_timeout,upstream_max_idle)
values
(162, 35, 0, 2000, 5000, 2, '127.0.0.1:50051', '100', '', 10000, 0, 0, 0),
(165, 34, 0, 2000, 5000, 2, '100.90.164.31:8072,100.90.163.51:8072,100.90.163.52:8072,100.90.165.32:8072', '50,50,50,80', '', 20000, 20000, 10000, 100),
(167, 36, 0, 2000, 5000, 2, '100.90.164.31:8072,100.90.163.51:8072,100.90.163.52:8072,100.90.165.32:8072', '50,50,50,80', '100.90.164.31:8072,100.90.163.51:8072', 10000, 10000, 10000, 100),
(168, 38, 0, 0, 0, 1, '111:111,22:111', '11,11', '111', 1111, 111, 222, 333),
(169, 41, 0, 0, 0, 1, '111:111,22:111', '11,11', '111', 0, 0, 0, 0),
(170, 42, 0, 0, 0, 1, '111:111,22:111', '11,11', '111', 0, 0, 0, 0),
(171, 43, 0, 2, 5, 1, '111:111,22:111', '11,11', '', 1111, 2222, 333, 444),
(172, 44, 0, 2, 5, 2, '127.0.0.1:8076', '50', '', 0, 0, 0, 0),
(173, 45, 0, 2, 5, 2, '127.0.0.1:88', '50', '', 0, 0, 0, 0),
(174, 46, 0, 2, 5, 2, '127.0.0.1:8002', '50', '', 0, 0, 0, 0),
(175, 47, 0, 2, 5, 2, '12777:11', '11', '', 0, 0, 0, 0),
(176, 48, 0, 2, 5, 2, '12777:11', '11', '', 0, 0, 0, 0),
(177, 49, 0, 2, 5, 2, '12777:11', '11', '', 0, 0, 0, 0),
(178, 50, 0, 2, 5, 2, '127.0.0.1:8001', '50', '', 0, 0, 0, 0),
(179, 51, 0, 2, 5, 2, '1212:11', '50', '', 0, 0, 0, 0),
(180, 52, 0, 2, 5, 2, '1212:11', '50', '', 0, 0, 0, 0),
(181, 53, 0, 2, 5, 2, '1111:11', '111', '', 0, 0, 0, 0),
(182, 54, 0, 2, 5, 1, '127.0.0.1:80', '50', '', 0, 0, 0, 0),
(183, 55, 0, 2, 5, 3, '127.0.0.1:81', '50', '', 0, 0, 0, 0),
(184, 56, 0, 2, 5, 2, '127.0.0.1:2003,127.0.0.1:2004', '50,50', '', 0, 0, 0, 0),
(185, 57, 0, 2, 5, 2, '127.0.0.1:6379', '50', '', 0, 0, 0, 0),
(186, 58, 0, 2, 5, 2, '127.0.0.1:50055', '50', '', 0, 0, 0, 0),
(187, 59, 0, 2, 5, 2, '127.0.0.1:2003,127.0.0.1:2004', '50,50', '', 0, 0, 0, 0),
(188, 60, 0, 2, 5, 2, '127.0.0.1:2003,127.0.0.1:2004', '50,50', '', 0, 0, 0, 0),
(189, 61, 0, 2, 5, 2, '127.0.0.1:3003,127.0.0.1:3004', '50,50', '', 0, 0, 0, 0);

/*
    select * from go_gateway.gateway_service_load_balance 查询测试数据是否添加成功
    delete from go_gateway.gateway_service_load_balance 删除测试数据
*/

-- -------------------------------------------------------------------
--
-- 表的结构 `go_gateway.gateway_service_tcp_rule`
--

create table if not exists go_gateway.gateway_service_tcp_rule (
    `id`            bigint(20)  not null                    comment '自增主键',
    `service_id`    bigint(20)  not null                    comment '服务id',
    `port`          int(5)      not null    default '0'     comment '端口号'
)engine=InnoDB default charset=utf8 comment='网关路由匹配表';

/* !drop table if exists go_gateway.gateway_service_tcp_rule */

--
--  向go_gateway.gateway_service_tcp_rule写入测试数据
--

insert into go_gateway.gateway_service_tcp_rule (id,service_id,port)
values
(171, 41, 8002),
(172, 42, 8003),
(173, 43, 8004),
(174, 38, 8004),
(175, 45, 8001),
(176, 46, 8005),
(177, 50, 8006),
(178, 51, 8007),
(179, 52, 8008),
(180, 55, 8010),
(181, 57, 8011);

/*
    select * from go_gateway.gateway_service_tcp_rule 查询测试数据是否添加成功
    delete from go_gateway.gateway_service_tcp_rule 删除测试数据
*/

-- ---------------------------------block end--------------------------------------|
-- --------------------------------------------------------------------------------|
-- --------------------------------------------------------------------------------|

-- --------------------------------------------------------------------------------|
-- --------------------------------------------------------------------------------|
-- ---------------------------------set indexes------------------------------------|

--
--  indexes for dumped tables
--

--
-- indexes for table gateway_admin
--
alter table gateway_admin
add primary key (`id`);

--
-- indexes for table gateway_app
--
alter table gateway_app
add primary key (`id`);

--
-- index for table gateway_service_access_control
--
alter table gateway_service_access_control
add primary key (`id`);

--
-- indexes for table gateway_service_grpc_rule
--
alter table gateway_service_grpc_rule
add primary key (`id`);

-- 
-- indexes for table gateway_service_http_rule
--
alter table gateway_service_http_rule
add primary key (`id`);

--
-- indexes for table gateway_service_info
--
alter table gateway_service_info
add primary key (`id`);

--
-- indexes for table gateway_service_load_balance
--
alter table gateway_service_load_balance
add primary key (`id`);

--
-- indexes for table gateway_service_tcp_rule
--
alter table gateway_service_tcp_rule
add primary key (`id`);

--
-- 在导出的表使用auto_increment
--

/* 
    auto_increment 指定一个列拥有自增属性
    具有auto_increment属性的数列应该是一个正数数列，如果把该数列声明为unsigned，这样序列的标号可增加一倍。比如tinyint数据列的最大编号是127，如果加上UNSIGNED，那么最大编号变为255
    auto_increment 数据列必须有唯一索引，以避免序号重复，必须具备not null属性
    如果把一个NULL插入到一个AUTO_INCREMENT数据列里去，MySQL将自动生成下一个序列编号。编号从1开始，并1为基数递增
    把0插入auto_increment数据列的效果与插入null值一样，但不建议这样做，而且在开头已经设置了SQL_MODE = `NO_AUTO_VALUE_ON_ZERO`，不允许插入0值
*/

--
-- 使用表auto_increment gateway_admin
--
alter table gateway_admin
modify `id` bigint(20) not null auto_increment comment '自增id', auto_increment=2;

--
-- 使用表auto_increment gateway_app
--
alter table gateway_app
modify `id` bigint(20) not null auto_increment comment '自增id', auto_increment=35;

--
-- 使用表auto_increment gateway_service_access_control
--
alter table gateway_service_access_control
modify `id` bigint(20) not null auto_increment comment '自增主键', auto_increment=190;

--
-- 使用表auto_increment gateway_service_grpc_rule
--
alter table gateway_service_grpc_rule
modify `id` bigint(20) not null auto_increment comment '自增主键', auto_increment=174;

--
-- 使用表auto_increment gateway_service_http_rule
--
alter table gateway_service_http_rule
modify `id` bigint(20) not null auto_increment comment '自增主键', auto_increment=181;

--
-- 使用表auto_increment gateway_service_info
--
alter table gateway_service_info
modify `id` bigint(20) unsigned not null auto_increment comment '自增主键', auto_increment=62;

--
-- 使用表auto_increment gateway_service_load_balance
--
alter table gateway_service_load_balance
modify `id` bigint(20) not null auto_increment comment '自增主键', auto_increment=190;

--
-- 使用表auto_increment gateway_service_tcp_rule
--
alter table gateway_service_tcp_rule
modify `id` bigint(20) not null auto_increment comment '自增主键', auto_increment=182;

-- ---------------------------------block end--------------------------------------|
-- --------------------------------------------------------------------------------|
-- --------------------------------------------------------------------------------|

/* 提交事务 */
commit;




