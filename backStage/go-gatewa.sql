/*
    Host:127.0.0.1
    Generation:2020-11-08
    server edition ： 5.7.32
*/

---------------------------------------------------------------------
/*
    创建数据库(go_dateway)
*/
create database if not exists go_gateway
default character set utf8
default collate utf8_general_ci;

/* !drop database if exists go_gateway*/
---------------------------------------------------------------------

/* 影响AUTO_INCREMETNT列的处理，
一般情况下你可以向该列插入NULL或0生成下一个序列号。
但是NO_AUTO_VALUE_ON_ZERO禁用0，
因此只有NULL可以生成下一个序列号 */
set SQL_MODE = "NO_AUTO_VALUE_ON_ZERO"; 

/* 当前session禁用自动提交事务，
自此句执行以后，每个SQL语句或者语句块所在的事务都需要显示commit才能提交事务 */
SET AUTOCOMMIT = 0; 

/* 启动一个新事务 */
START TRANSACTION; 
/* commit */

/* 修改当前会话时区 */
SET time_zone="+00:00"; 


---------------------------------------------------------------------
-- 
-- Database : "go_gateway"
--

--

--
-- 表的结构："gateway_admin"
--

CREATE TABLE if not exists go_gateway.gateway_admin  (
    `id` bigint(20) not null comment '自增id',
    `user_name` varchar(255) not null default '' comment '用户名',
    `salt` varchar(50) not null default '' comment '盐',
    `password` varchar(255) not null default '' comment '密码',
    `create_at` datetime not null default '1971-01-01 00:00:00' comment '新增时间' ,
    `update_at` datetime not null default '1971-01-01 00:00:00' comment '更新时间',
    `is_delete` tinyint(4) not null default '0' comment '是否删除'
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

---------------------------------------------------------------------
--
--
--
