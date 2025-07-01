#!/bin/bash

# 加载环境变量
export $(cat .env | xargs)

# 运行程序
go run cmd/dns-update/main.go 