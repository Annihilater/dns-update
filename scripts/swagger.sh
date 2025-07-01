#!/bin/bash

# 确保swag命令存在
if ! command -v swag &> /dev/null; then
    echo "Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# 生成swagger文档
echo "Generating swagger docs..."
swag init -g cmd/dns-update/main.go

echo "Swagger docs updated successfully!"
echo "You can view the API documentation at: http://localhost:8080/swagger/index.html" 