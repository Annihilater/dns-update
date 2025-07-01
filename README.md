# DNS Update

阿里云 DNS 解析管理工具，支持域名和解析记录的增删改查操作。

## 功能特性

- 域名管理
  - 查询账户下域名列表
  - 添加新域名
  - 查询域名信息
  
- 解析记录管理
  - 添加解析记录
  - 更新解析记录
  - 删除解析记录
  - 查询解析记录
  - 设置解析状态
  
- 域名分组管理
  - 查询域名组
  - 添加域名组
  - 更新域名组
  - 删除域名组

## 环境要求

- Go 1.16 或更高版本
- 阿里云账号和 AccessKey

## 安装

```bash
go get github.com/yourusername/dns-update
```

## 配置

1. 设置环境变量：

```bash
export ACCESS_KEY_ID=your_access_key_id
export ACCESS_KEY_SECRET=your_access_key_secret
```

2. 或者修改 `configs/config.yaml`：

```yaml
aliyun:
  access_key_id: your_access_key_id
  access_key_secret: your_access_key_secret
  region_id: cn-hangzhou
```

## 使用方法

```bash
dns-update <domainName> <RR> <recordType> <value> <recordId> <groupName> <groupId> <regionId>
```

参数说明：
- domainName: 域名
- RR: 主机记录
- recordType: 记录类型(A/NS/MX/TXT/CNAME/SRV/AAAA/CAA/REDIRECT_URL/FORWARD_URL)
- value: 记录值
- recordId: 解析记录ID
- groupName: 域名组名称
- groupId: 域名组ID
- regionId: 地域ID

## 项目结构

```
.
├── cmd/                    # 主要的应用程序入口
│   └── dns-update/        # 应用程序名称
│       └── main.go        # 主程序入口文件
├── internal/              # 私有应用程序和库代码
│   ├── handler/          # 处理程序
│   ├── service/          # 业务逻辑
│   └── repository/       # 数据访问层
├── pkg/                   # 可以被外部应用程序使用的库代码
├── configs/              # 配置文件
├── docs/                 # 文档
└── README.md             # 项目说明
```

## 开发

1. 克隆仓库：

```bash
git clone https://github.com/yourusername/dns-update.git
```

2. 安装依赖：

```bash
go mod tidy
```

3. 构建：

```bash
go build -o dns-update ./cmd/dns-update
```

## 贡献

欢迎提交 Issue 和 Pull Request。

## 许可证

MIT License 