#!/bin/sh

# 开启GO111MODULE
export GO111MODULE=on
# 设置代理服务器
export GOPROXY=https://goproxy.io
# 生成bin文件
go build -o bin/gatewayDemo
# 检查当前晋城市会否有gatewayDemo，如果有就杀死进程
ps aux | grep gatewayDemo | grep -v 'grep' | awk '{print $2}' | xargs kill -9
#
nohup ./bin/gatewayDemo -conf=./conf/prod/ -endpoint=dashboard >>logs/dashboard.log 2>&1 &
echo 'nohup ./bin/gatewayDemo -conf=./conf/prod/ -endpoint=dashboard >> logs/dashboard.log 2>&1 &'
nohup ./bin/gatewayDemo -conf=./conf/prod/ -endpoint=server >>logs/server.log 2>&1 &
echo "nohup ./bin/gatewayDemo -conf=./conf/prod/ -endpoint=server >> logs/server.log 2>&1 &"
